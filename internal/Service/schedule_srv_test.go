package Service

import (
	"task/internal/Service/mocks"

	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task/internal/Models"
)

func TestScheduleService_CreateSchedule_Success(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	svc := NewScheduleService(mockScheduleRepo, mockSlotRepo, mockRoomRepo)
	roomID := uuid.New()
	schedule := &Models.Schedule{
		ID:         uuid.New(),
		RoomID:     roomID,
		DaysOfWeek: []int{1, 2, 3, 4, 5},
		StartTime:  "09:00",
		EndTime:    "17:00",
	}
	mockRoomRepo.On("Exists", roomID).Return(true, nil)
	mockScheduleRepo.On("ExistsForRoom", roomID).Return(false, nil)
	mockScheduleRepo.On("Create", mock.AnythingOfType("*Models.Schedule")).Return(nil)
	mockSlotRepo.On("DeleteByRoomID", roomID).Return(nil)
	mockSlotRepo.On("CreateBatch", mock.Anything).Return(nil)
	err := svc.CreateSchedule(schedule)
	assert.NoError(t, err)
	mockRoomRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
	mockSlotRepo.AssertExpectations(t)
}

func TestScheduleService_CreateSchedule_RoomNotFound(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	svc := NewScheduleService(mockScheduleRepo, mockSlotRepo, mockRoomRepo)
	roomID := uuid.New()
	schedule := &Models.Schedule{RoomID: roomID}
	mockRoomRepo.On("Exists", roomID).Return(false, nil)
	err := svc.CreateSchedule(schedule)
	assert.ErrorIs(t, err, ErrRoomNotFound)
}

func TestScheduleService_CreateSchedule_AlreadyExists(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	svc := NewScheduleService(mockScheduleRepo, mockSlotRepo, mockRoomRepo)
	roomID := uuid.New()
	schedule := &Models.Schedule{RoomID: roomID}
	mockRoomRepo.On("Exists", roomID).Return(true, nil)
	mockScheduleRepo.On("ExistsForRoom", roomID).Return(true, nil)
	err := svc.CreateSchedule(schedule)
	assert.ErrorIs(t, err, ErrScheduleAlreadyExists)
}

func TestScheduleService_CreateSchedule_InvalidTime(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	svc := NewScheduleService(mockScheduleRepo, mockSlotRepo, mockRoomRepo)
	roomID := uuid.New()
	schedule := &Models.Schedule{
		RoomID:     roomID,
		StartTime:  "17:00",
		EndTime:    "09:00",
		DaysOfWeek: []int{1},
	}

	mockRoomRepo.On("Exists", roomID).Return(true, nil)
	mockScheduleRepo.On("ExistsForRoom", roomID).Return(false, nil)
	err := svc.CreateSchedule(schedule)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start time must be before end time")
}
