package mocks

import (
	"time"

	"task/internal/Models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) Create(room *Models.Room) error {
	args := m.Called(room)
	return args.Error(0)
}
func (m *MockRoomRepository) List() ([]Models.Room, error) {
	args := m.Called()
	return args.Get(0).([]Models.Room), args.Error(1)
}
func (m *MockRoomRepository) Exists(id uuid.UUID) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

type MockScheduleRepository struct {
	mock.Mock
}

func (m *MockScheduleRepository) Create(schedule *Models.Schedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}
func (m *MockScheduleRepository) FindByRoomID(roomID uuid.UUID) (*Models.Schedule, error) {
	args := m.Called(roomID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Models.Schedule), args.Error(1)
}
func (m *MockScheduleRepository) ExistsForRoom(roomID uuid.UUID) (bool, error) {
	args := m.Called(roomID)
	return args.Bool(0), args.Error(1)
}

type MockSlotRepository struct {
	mock.Mock
}

func (m *MockSlotRepository) Create(slot *Models.Slot) error {
	args := m.Called(slot)
	return args.Error(0)
}
func (m *MockSlotRepository) CreateBatch(slots []Models.Slot) error {
	args := m.Called(slots)
	return args.Error(0)
}
func (m *MockSlotRepository) FindAvailableByRoomAndDate(roomID uuid.UUID, date time.Time) ([]Models.Slot, error) {
	args := m.Called(roomID, date)
	return args.Get(0).([]Models.Slot), args.Error(1)
}
func (m *MockSlotRepository) FindByID(id uuid.UUID) (*Models.Slot, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Models.Slot), args.Error(1)
}
func (m *MockSlotRepository) DeleteByRoomID(roomID uuid.UUID) error {
	args := m.Called(roomID)
	return args.Error(0)
}

type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) Create(booking *Models.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}
func (m *MockBookingRepository) Cancel(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockBookingRepository) FindByID(id uuid.UUID) (*Models.Booking, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Models.Booking), args.Error(1)
}
func (m *MockBookingRepository) FindActiveBySlotID(slotID uuid.UUID) (*Models.Booking, error) {
	args := m.Called(slotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Models.Booking), args.Error(1)
}
func (m *MockBookingRepository) FindByUserID(userID uuid.UUID, onlyFuture bool) ([]Models.Booking, error) {
	args := m.Called(userID, onlyFuture)
	return args.Get(0).([]Models.Booking), args.Error(1)
}
func (m *MockBookingRepository) ListAll(page, pageSize int) ([]Models.Booking, int, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]Models.Booking), args.Int(1), args.Error(2)
}
