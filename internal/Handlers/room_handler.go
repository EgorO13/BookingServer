package Handlers

import (
	"encoding/json"
	"net/http"

	"task/internal/Service"
)

type RoomHandler struct {
	IRoomService Service.IRoomService
}

func NewRoomHandler(RoomService Service.IRoomService) *RoomHandler {
	return &RoomHandler{IRoomService: RoomService}
}

type createRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Capacity    int    `json:"capacity,omitempty"`
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req createRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		SendErrorResponse(w, "INVALID_REQUEST", "name is required", http.StatusBadRequest)
		return
	}
	var description *string
	if req.Description != "" {
		description = &req.Description
	}
	var capacity *int
	if req.Capacity != 0 {
		capacity = &req.Capacity
	}
	room, err := h.IRoomService.CreateRoom(req.Name, description, capacity)
	if err != nil {
		SendErrorResponse(w, "INTERNAL_ERROR", "failed to create room", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"room": room})
}

func (h *RoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.IRoomService.ListRooms()
	if err != nil {
		SendErrorResponse(w, "INTERNAL_ERROR", "failed to list rooms", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"rooms": rooms})
}
