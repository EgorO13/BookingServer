package main

import (
	"log"

	"task/internal/Config"
	"task/internal/DB"
	"task/internal/Models"
	"task/internal/Repository"
	"task/internal/Service"

	"github.com/google/uuid"
)

func main() {
	cfg := config.Load()
	database, err := db.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer database.Close()
	roomRepo := Repository.NewRoomRepository(database)
	scheduleRepo := Repository.NewScheduleRepository(database)
	slotRepo := Repository.NewSlotRepository(database)
	room := &Models.Room{
		ID:          uuid.New(),
		Name:        "Test room",
		Description: stringPtr("Auto"),
		Capacity:    intPtr(5),
	}
	if err := roomRepo.Create(room); err != nil {
		log.Fatal("Failed to create room:", err)
	}
	log.Printf("Room created: %s", room.ID)
	scheduleService := Service.NewScheduleService(scheduleRepo, slotRepo, roomRepo)
	schedule := &Models.Schedule{
		ID:         uuid.New(),
		RoomID:     room.ID,
		DaysOfWeek: []int{1, 2, 3, 4, 5},
		StartTime:  "09:00",
		EndTime:    "18:00",
	}
	if err := scheduleService.CreateSchedule(schedule); err != nil {
		log.Fatal("Failed to create schedule and generate slots:", err)
	}
	log.Printf("Schedule created and slots generated for room %s", room.ID)
}

func stringPtr(s string) *string { return &s }
func intPtr(i int) *int          { return &i }
