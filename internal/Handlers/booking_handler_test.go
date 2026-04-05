package Handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"task/internal/Auth"
	"task/internal/Models"
)

type mockBookingService struct {
	mock.Mock
}

func (m *mockBookingService) CreateBooking(userID, slotID uuid.UUID) (*Models.Booking, error) {
	args := m.Called(userID, slotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Models.Booking), args.Error(1)
}
func (m *mockBookingService) CancelBooking(bookingID, userID uuid.UUID) (*Models.Booking, error) {
	args := m.Called(bookingID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Models.Booking), args.Error(1)
}
func (m *mockBookingService) GetMyBookings(userID uuid.UUID) ([]Models.Booking, error) {
	args := m.Called(userID)
	return args.Get(0).([]Models.Booking), args.Error(1)
}
func (m *mockBookingService) GetAllBookings(page, pageSize int) ([]Models.Booking, int, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]Models.Booking), args.Int(1), args.Error(2)
}

func TestBookingHandler_CreateBooking(t *testing.T) {
	mockSvc := new(mockBookingService)
	handler := NewBookingHandler(mockSvc)
	slotID := uuid.New()
	reqBody := `{"slotId":"` + slotID.String() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewReader([]byte(reqBody)))
	ctx := context.WithValue(req.Context(), Auth.UserIDKey, "00000000-0000-0000-0000-000000000002")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	expectedBooking := &Models.Booking{ID: uuid.New(), SlotID: slotID, Status: "active"}
	mockSvc.On("CreateBooking", uuid.MustParse("00000000-0000-0000-0000-000000000002"), slotID).Return(expectedBooking, nil)
	handler.CreateBooking(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.NotNil(t, response["booking"])
	mockSvc.AssertExpectations(t)
}
