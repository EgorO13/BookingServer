package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"task/internal/Auth"
	"task/internal/Config"
	"task/internal/DB"
	"task/internal/Handlers"
	"task/internal/Repository"
	"task/internal/Service"
)

func main() {
	cfg := config.Load()

	database, err := db.NewDB(cfg)
	if err != nil {
		log.Fatal("DB connect error: ", err)
	}
	defer database.Close()

	migrator, err := db.NewMigrationDriver(database)
	if err != nil {
		log.Fatal("Failed to create migration driver: ", err)
	}
	if err := db.RunMigrations(migrator); err != nil {
		log.Fatal("Failed to apply migrations: ", err)
	}

	roomRepo := Repository.NewRoomRepository(database)
	scheduleRepo := Repository.NewScheduleRepository(database)
	slotRepo := Repository.NewSlotRepository(database)
	bookingRepo := Repository.NewBookingRepository(database)

	roomService := Service.NewRoomService(roomRepo)
	scheduleService := Service.NewScheduleService(scheduleRepo, slotRepo, roomRepo)
	slotService := Service.NewSlotService(slotRepo, bookingRepo, roomRepo)
	bookingService := Service.NewBookingService(bookingRepo, slotRepo, database)

	roomHandler := Handlers.NewRoomHandler(roomService)
	scheduleHandler := Handlers.NewScheduleHandler(scheduleService)
	slotHandler := Handlers.NewSlotHandler(slotService)
	bookingHandler := Handlers.NewBookingHandler(bookingService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Post("/dummyLogin", Handlers.DummyLoginHandler(cfg.JWTSecret))

	r.Route("/rooms", func(r chi.Router) {
		r.With(Auth.AuthMiddleware(cfg.JWTSecret)).Get("/list", roomHandler.ListRooms)
		r.With(Auth.AuthMiddleware(cfg.JWTSecret)).Post("/create", Auth.AdminOnly(roomHandler.CreateRoom))
	})

	r.Route("/rooms/{roomId}/schedule", func(r chi.Router) {
		r.With(Auth.AuthMiddleware(cfg.JWTSecret)).Post("/create", Auth.AdminOnly(scheduleHandler.CreateSchedule))
	})

	r.Route("/rooms/{roomId}/slots", func(r chi.Router) {
		r.With(Auth.AuthMiddleware(cfg.JWTSecret)).Get("/list", slotHandler.ListAvailableSlots)
	})

	r.Route("/bookings", func(r chi.Router) {
		r.With(Auth.AuthMiddleware(cfg.JWTSecret)).Post("/create", Auth.UserOnly(bookingHandler.CreateBooking))
		r.With(Auth.AuthMiddleware(cfg.JWTSecret)).Post("/{bookingId}/cancel", Auth.UserOnly(bookingHandler.CancelBooking))
		r.With(Auth.AuthMiddleware(cfg.JWTSecret)).Get("/my", Auth.UserOnly(bookingHandler.MyBookings))
		r.With(Auth.AuthMiddleware(cfg.JWTSecret)).Get("/list", Auth.AdminOnly(bookingHandler.ListAllBookings))
	})

	port := ":8080"
	fmt.Printf("Server starting on %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
