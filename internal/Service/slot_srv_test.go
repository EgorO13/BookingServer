package Service

import (
	"task/internal/Service/mocks"
	"testing"
	"time"

	"task/internal/Models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSlotService_GetAvailableSlots(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	mockBookingRepo := new(mocks.MockBookingRepository)
	svc := NewSlotService(mockSlotRepo, mockBookingRepo, mockRoomRepo)
	roomID := uuid.New()
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedSlots := []Models.Slot{{ID: uuid.New(), RoomID: roomID}}
	mockRoomRepo.On("Exists", roomID).Return(true, nil)
	mockSlotRepo.On("FindAvailableByRoomAndDate", roomID, date).Return(expectedSlots, nil)
	slots, err := svc.GetAvailableSlots(roomID, date)
	assert.NoError(t, err)
	assert.Equal(t, expectedSlots, slots)
	mockRoomRepo.AssertExpectations(t)
	mockSlotRepo.AssertExpectations(t)
}

func TestSlotService_GetAvailableSlots_RoomNotFound(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	mockBookingRepo := new(mocks.MockBookingRepository)
	svc := NewSlotService(mockSlotRepo, mockBookingRepo, mockRoomRepo)
	roomID := uuid.New()
	date := time.Now()
	mockRoomRepo.On("Exists", roomID).Return(false, nil)
	slots, err := svc.GetAvailableSlots(roomID, date)
	assert.ErrorIs(t, err, ErrRoomNotFound)
	assert.Nil(t, slots)
}

func TestSlotService_GetSlotByID(t *testing.T) {
	mockSlotRepo := new(mocks.MockSlotRepository)
	mockBookingRepo := new(mocks.MockBookingRepository)
	mockRoomRepo := new(mocks.MockRoomRepository)
	svc := NewSlotService(mockSlotRepo, mockBookingRepo, mockRoomRepo)
	slotID := uuid.New()
	expected := &Models.Slot{ID: slotID}
	mockSlotRepo.On("FindByID", slotID).Return(expected, nil)
	slot, err := svc.GetSlotByID(slotID)
	assert.NoError(t, err)
	assert.Equal(t, expected, slot)
}

func TestSlotService_GetSlotByID_NotFound(t *testing.T) {
	mockSlotRepo := new(mocks.MockSlotRepository)
	mockBookingRepo := new(mocks.MockBookingRepository)
	mockRoomRepo := new(mocks.MockRoomRepository)
	svc := NewSlotService(mockSlotRepo, mockBookingRepo, mockRoomRepo)
	slotID := uuid.New()
	mockSlotRepo.On("FindByID", slotID).Return(nil, nil)
	slot, err := svc.GetSlotByID(slotID)
	assert.ErrorIs(t, err, ErrSlotNotFound)
	assert.Nil(t, slot)
}
