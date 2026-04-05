package Service

import (
	"fmt"
	"time"

	"task/internal/Models"
	"task/internal/Repository"

	"github.com/google/uuid"
)

type ISlotService interface {
	GetAvailableSlots(roomID uuid.UUID, date time.Time) ([]Models.Slot, error)
	GetSlotByID(id uuid.UUID) (*Models.Slot, error)
}

type SlotService struct {
	slotRepo    Repository.ISlotRepository
	bookingRepo Repository.IBookingRepository
	roomRepo    Repository.IRoomRepository
}

func NewSlotService(slotRepo Repository.ISlotRepository, bookingRepo Repository.IBookingRepository, roomRepo Repository.IRoomRepository) ISlotService {
	return &SlotService{
		slotRepo:    slotRepo,
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
	}
}

func (s *SlotService) GetAvailableSlots(roomID uuid.UUID, date time.Time) ([]Models.Slot, error) {
	exists, err := s.roomRepo.Exists(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to check room: %w", err)
	}
	if !exists {
		return nil, ErrRoomNotFound
	}
	dateUTC := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	slots, err := s.slotRepo.FindAvailableByRoomAndDate(roomID, dateUTC)
	if err != nil {
		return nil, fmt.Errorf("failed to get available slots: %w", err)
	}
	return slots, nil
}

func (s *SlotService) GetSlotByID(id uuid.UUID) (*Models.Slot, error) {
	slot, err := s.slotRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find slot: %w", err)
	}
	if slot == nil {
		return nil, ErrSlotNotFound
	}
	return slot, nil
}
