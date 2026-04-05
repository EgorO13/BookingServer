package Handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"task/internal/Models"
	"task/internal/Service"
)

type ScheduleHandler struct {
	ScheduleService Service.IScheduleService
}

func NewScheduleHandler(scheduleService Service.IScheduleService) *ScheduleHandler {
	return &ScheduleHandler{ScheduleService: scheduleService}
}

type createScheduleRequest struct {
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		SendErrorResponse(w, "INVALID_REQUEST", "invalid roomId", http.StatusBadRequest)
		return
	}

	var req createScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}

	for _, day := range req.DaysOfWeek {
		if day < 1 || day > 7 {
			SendErrorResponse(w, "INVALID_REQUEST", "daysOfWeek must be between 1 and 7", http.StatusBadRequest)
			return
		}
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		SendErrorResponse(w, "INVALID_REQUEST", "invalid startTime format, use HH:MM", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		SendErrorResponse(w, "INVALID_REQUEST", "invalid endTime format, use HH:MM", http.StatusBadRequest)
		return
	}
	if !startTime.Before(endTime) {
		SendErrorResponse(w, "INVALID_REQUEST", "startTime must be before endTime", http.StatusBadRequest)
		return
	}

	schedule := &Models.Schedule{
		ID:         uuid.New(),
		RoomID:     roomID,
		DaysOfWeek: req.DaysOfWeek,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}

	if err := h.ScheduleService.CreateSchedule(schedule); err != nil {
		switch err {
		case Service.ErrRoomNotFound:
			SendErrorResponse(w, "ROOM_NOT_FOUND", "room not found", http.StatusNotFound)
		case Service.ErrScheduleAlreadyExists:
			SendErrorResponse(w, "SCHEDULE_EXISTS", "schedule already exists for this room", http.StatusConflict)
		default:
			SendErrorResponse(w, "INTERNAL_ERROR", "failed to create schedule", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"schedule": schedule})
}
