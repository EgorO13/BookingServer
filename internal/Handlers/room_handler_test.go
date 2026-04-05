package Handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"task/internal/Models"
)

type mockRoomService struct {
	mock.Mock
}

func (m *mockRoomService) CreateRoom(name string, description *string, capacity *int) (*Models.Room, error) {
	args := m.Called(name, description, capacity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Models.Room), args.Error(1)
}
func (m *mockRoomService) ListRooms() ([]Models.Room, error) {
	args := m.Called()
	return args.Get(0).([]Models.Room), args.Error(1)
}

func TestRoomHandler_CreateRoom(t *testing.T) {
	mockSvc := new(mockRoomService)
	handler := NewRoomHandler(mockSvc)
	reqBody := `{"name":"Test Room","description":"desc","capacity":10}`
	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewReader([]byte(reqBody)))
	w := httptest.NewRecorder()
	expectedRoom := &Models.Room{ID: uuid.New(), Name: "Test Room"}
	mockSvc.On("CreateRoom", "Test Room", stringPtr("desc"), intPtr(10)).Return(expectedRoom, nil)
	handler.CreateRoom(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.NotNil(t, response["room"])
	mockSvc.AssertExpectations(t)
}

func TestRoomHandler_ListRooms(t *testing.T) {
	mockSvc := new(mockRoomService)
	handler := NewRoomHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
	w := httptest.NewRecorder()
	expectedRooms := []Models.Room{{ID: uuid.New(), Name: "Room1"}, {ID: uuid.New(), Name: "Room2"}}
	mockSvc.On("ListRooms").Return(expectedRooms, nil)
	handler.ListRooms(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	rooms := response["rooms"].([]interface{})
	assert.Len(t, rooms, 2)
	mockSvc.AssertExpectations(t)
}

func stringPtr(s string) *string { return &s }
func intPtr(i int) *int          { return &i }
