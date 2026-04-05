package Handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task/internal/Models"
)

type mockScheduleService struct {
	mock.Mock
}

func (m *mockScheduleService) CreateSchedule(schedule *Models.Schedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}

func TestScheduleHandler_CreateSchedule(t *testing.T) {
	mockSvc := new(mockScheduleService)
	handler := NewScheduleHandler(mockSvc)
	roomID := uuid.New()
	scheduleBody := `{"daysOfWeek":[1,2,3,4,5],"startTime":"09:00","endTime":"17:00"}`
	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID.String()+"/schedule/create", bytes.NewReader([]byte(scheduleBody)))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("roomId", roomID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	mockSvc.On("CreateSchedule", mock.AnythingOfType("*Models.Schedule")).Return(nil)
	handler.CreateSchedule(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockSvc.AssertExpectations(t)
}
