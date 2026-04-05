package Handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task/internal/Auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"task/internal/Service"
)

type BookingHandler struct {
	BookingService Service.IBookingService
}

func NewBookingHandler(BookingService Service.IBookingService) *BookingHandler {
	return &BookingHandler{BookingService: BookingService}
}

type createBookingRequest struct {
	SlotID string `json:"slotId"`
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(Auth.UserIDKey).(string)
	if !ok {
		SendErrorResponse(w, "UNAUTHORIZED", "user not authenticated", http.StatusUnauthorized)
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		SendErrorResponse(w, "UNAUTHORIZED", "invalid user id", http.StatusUnauthorized)
		return
	}

	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}
	slotID, err := uuid.Parse(req.SlotID)
	if err != nil {
		SendErrorResponse(w, "INVALID_REQUEST", "invalid slotId", http.StatusBadRequest)
		return
	}

	booking, err := h.BookingService.CreateBooking(userUUID, slotID)
	if err != nil {
		switch err {
		case Service.ErrSlotNotFound:
			SendErrorResponse(w, "SLOT_NOT_FOUND", "slot not found", http.StatusNotFound)
		case Service.ErrSlotInPast:
			SendErrorResponse(w, "INVALID_REQUEST", "cannot book slot in the past", http.StatusBadRequest)
		case Service.ErrSlotAlreadyBooked:
			SendErrorResponse(w, "SLOT_ALREADY_BOOKED", "slot is already booked", http.StatusConflict)
		default:
			SendErrorResponse(w, "INTERNAL_ERROR", "failed to create booking", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"booking": booking})
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(Auth.UserIDKey).(string)
	if !ok {
		SendErrorResponse(w, "UNAUTHORIZED", "user not authenticated", http.StatusUnauthorized)
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		SendErrorResponse(w, "UNAUTHORIZED", "invalid user id", http.StatusUnauthorized)
		return
	}

	bookingIDStr := chi.URLParam(r, "bookingId")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		SendErrorResponse(w, "INVALID_REQUEST", "invalid bookingId", http.StatusBadRequest)
		return
	}

	booking, err := h.BookingService.CancelBooking(bookingID, userUUID)
	if err != nil {
		switch err {
		case Service.ErrBookingNotFound:
			SendErrorResponse(w, "BOOKING_NOT_FOUND", "booking not found", http.StatusNotFound)
		case Service.ErrNotOwner:
			SendErrorResponse(w, "FORBIDDEN", "you can only cancel your own bookings", http.StatusForbidden)
		default:
			SendErrorResponse(w, "INTERNAL_ERROR", "failed to cancel booking", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"booking": booking})
}

func (h *BookingHandler) MyBookings(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(Auth.UserIDKey).(string)
	if !ok {
		SendErrorResponse(w, "UNAUTHORIZED", "user not authenticated", http.StatusUnauthorized)
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		SendErrorResponse(w, "UNAUTHORIZED", "invalid user id", http.StatusUnauthorized)
		return
	}

	bookings, err := h.BookingService.GetMyBookings(userUUID)
	if err != nil {
		SendErrorResponse(w, "INTERNAL_ERROR", "failed to get bookings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"bookings": bookings})
}

func (h *BookingHandler) ListAllBookings(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 1 {
			page = p
		}
	}
	pageSize := 20
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps >= 1 && ps <= 100 {
			pageSize = ps
		}
	}

	bookings, total, err := h.BookingService.GetAllBookings(page, pageSize)
	if err != nil {
		SendErrorResponse(w, "INTERNAL_ERROR", "failed to list bookings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"bookings": bookings,
		"pagination": map[string]interface{}{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	})
}
