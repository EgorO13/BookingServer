package Handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task/internal/Models"
)

type mockSlotService struct {
	mock.Mock
}

func (m *mockSlotService) GetAvailableSlots(roomID uuid.UUID, date time.Time) ([]Models.Slot, error) {
	args := m.Called(roomID, date)
	return args.Get(0).([]Models.Slot), args.Error(1)
}
func (m *mockSlotService) GetSlotByID(id uuid.UUID) (*Models.Slot, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Models.Slot), args.Error(1)
}

func TestSlotHandler_ListAvailableSlots(t *testing.T) {
	mockSvc := new(mockSlotService)
	handler := NewSlotHandler(mockSvc)
	roomID := uuid.New()
	date := "2024-01-15"
	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/slots/list?date="+date, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("roomId", roomID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	expectedSlots := []Models.Slot{{ID: uuid.New(), RoomID: roomID}}
	mockSvc.On("GetAvailableSlots", roomID, mock.AnythingOfType("time.Time")).Return(expectedSlots, nil)
	handler.ListAvailableSlots(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockSvc.AssertExpectations(t)
}
