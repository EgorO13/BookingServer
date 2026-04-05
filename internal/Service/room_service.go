package Service

import (
	"fmt"

	"task/internal/Models"
	"task/internal/Repository"

	"github.com/google/uuid"
)

type IRoomService interface {
	CreateRoom(name string, description *string, capacity *int) (*Models.Room, error)
	ListRooms() ([]Models.Room, error)
}

type RoomService struct {
	roomRepo Repository.IRoomRepository
}

func NewRoomService(roomRepo Repository.IRoomRepository) IRoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) CreateRoom(name string, description *string, capacity *int) (*Models.Room, error) {
	room := &Models.Room{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Capacity:    capacity,
	}
	if err := s.roomRepo.Create(room); err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}
	return room, nil
}

func (s *RoomService) ListRooms() ([]Models.Room, error) {
	rooms, err := s.roomRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list rooms: %w", err)
	}
	return rooms, nil
}
