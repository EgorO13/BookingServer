package Service

import (
	"task/internal/Service/mocks"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"task/internal/Models"
)

func TestBookingService_GetMyBookings(t *testing.T) {
	mockBookingRepo := new(mocks.MockBookingRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	svc := NewBookingService(mockBookingRepo, mockSlotRepo, db)
	userID := uuid.New()
	expected := []Models.Booking{{ID: uuid.New(), UserID: userID}}
	mockBookingRepo.On("FindByUserID", userID, true).Return(expected, nil)
	bookings, err := svc.GetMyBookings(userID)
	assert.NoError(t, err)
	assert.Equal(t, expected, bookings)
	mockBookingRepo.AssertExpectations(t)
}

func TestBookingService_GetAllBookings(t *testing.T) {
	mockBookingRepo := new(mocks.MockBookingRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	svc := NewBookingService(mockBookingRepo, mockSlotRepo, db)
	expected := []Models.Booking{{ID: uuid.New()}}
	total := 1
	mockBookingRepo.On("ListAll", 1, 20).Return(expected, total, nil)
	bookings, count, err := svc.GetAllBookings(1, 20)
	assert.NoError(t, err)
	assert.Equal(t, expected, bookings)
	assert.Equal(t, total, count)
}

func TestBookingService_CancelBooking_Own(t *testing.T) {
	mockBookingRepo := new(mocks.MockBookingRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	svc := NewBookingService(mockBookingRepo, mockSlotRepo, db)
	bookingID := uuid.New()
	userID := uuid.New()
	booking := &Models.Booking{ID: bookingID, UserID: userID, Status: string(Models.BookingStatusActive)}
	mockBookingRepo.On("FindByID", bookingID).Return(booking, nil)
	mockBookingRepo.On("Cancel", bookingID).Return(nil)
	result, err := svc.CancelBooking(bookingID, userID)
	assert.NoError(t, err)
	assert.Equal(t, string(Models.BookingStatusCancelled), result.Status)
	mockBookingRepo.AssertExpectations(t)
}

func TestBookingService_CancelBooking_AlreadyCancelled(t *testing.T) {
	mockBookingRepo := new(mocks.MockBookingRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	svc := NewBookingService(mockBookingRepo, mockSlotRepo, db)
	bookingID := uuid.New()
	userID := uuid.New()
	booking := &Models.Booking{ID: bookingID, UserID: userID, Status: string(Models.BookingStatusCancelled)}
	mockBookingRepo.On("FindByID", bookingID).Return(booking, nil)
	result, err := svc.CancelBooking(bookingID, userID)
	assert.NoError(t, err)
	assert.Equal(t, string(Models.BookingStatusCancelled), result.Status)
	mockBookingRepo.AssertNotCalled(t, "Cancel")
}

func TestBookingService_CancelBooking_NotOwner(t *testing.T) {
	mockBookingRepo := new(mocks.MockBookingRepository)
	mockSlotRepo := new(mocks.MockSlotRepository)
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	svc := NewBookingService(mockBookingRepo, mockSlotRepo, db)
	bookingID := uuid.New()
	ownerID := uuid.New()
	otherUserID := uuid.New()
	booking := &Models.Booking{ID: bookingID, UserID: ownerID}
	mockBookingRepo.On("FindByID", bookingID).Return(booking, nil)
	_, err = svc.CancelBooking(bookingID, otherUserID)
	assert.ErrorIs(t, err, ErrNotOwner)
}
