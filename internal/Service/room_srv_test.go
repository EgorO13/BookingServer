package Service

import (
	"task/internal/Service/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task/internal/Models"
)

func TestRoomService_CreateRoom(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	svc := NewRoomService(mockRepo)
	name := "Conference"
	desc := "Big room"
	cap := 20
	mockRepo.On("Create", mock.AnythingOfType("*Models.Room")).Return(nil).Once()
	room, err := svc.CreateRoom(name, &desc, &cap)
	assert.NoError(t, err)
	assert.Equal(t, name, room.Name)
	assert.Equal(t, desc, *room.Description)
	assert.Equal(t, cap, *room.Capacity)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_ListRooms(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	svc := NewRoomService(mockRepo)
	expected := []Models.Room{{ID: uuid.New(), Name: "A"}, {ID: uuid.New(), Name: "B"}}
	mockRepo.On("List").Return(expected, nil).Once()
	rooms, err := svc.ListRooms()
	assert.NoError(t, err)
	assert.Equal(t, expected, rooms)
	mockRepo.AssertExpectations(t)
}
