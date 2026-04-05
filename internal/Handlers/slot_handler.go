package Handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"task/internal/Service"
)

type SlotHandler struct {
	SlotService Service.ISlotService
}

func NewSlotHandler(SlotService Service.ISlotService) *SlotHandler {
	return &SlotHandler{SlotService: SlotService}
}

func (h *SlotHandler) ListAvailableSlots(w http.ResponseWriter, r *http.Request) {
	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		SendErrorResponse(w, "INVALID_REQUEST", "invalid roomId", http.StatusBadRequest)
		return
	}
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		SendErrorResponse(w, "INVALID_REQUEST", "date parameter is required", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		SendErrorResponse(w, "INVALID_REQUEST", "invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	slots, err := h.SlotService.GetAvailableSlots(roomID, date)
	if err != nil {
		switch err {
		case Service.ErrRoomNotFound:
			SendErrorResponse(w, "ROOM_NOT_FOUND", "room not found", http.StatusNotFound)
		default:
			SendErrorResponse(w, "INTERNAL_ERROR", "failed to get slots", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"slots": slots})
}
