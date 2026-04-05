package Service

import (
	"database/sql"
	"fmt"
	"time"

	"task/internal/Models"
	"task/internal/Repository"

	"github.com/google/uuid"
)

type IBookingService interface {
	CreateBooking(userID, slotID uuid.UUID) (*Models.Booking, error)
	CancelBooking(bookingID, userID uuid.UUID) (*Models.Booking, error)
	GetMyBookings(userID uuid.UUID) ([]Models.Booking, error)
	GetAllBookings(page, pageSize int) ([]Models.Booking, int, error)
}

type BookingService struct {
	bookingRepo Repository.IBookingRepository
	slotRepo    Repository.ISlotRepository
	db          *sql.DB
}

func NewBookingService(bookingRepo Repository.IBookingRepository, slotRepo Repository.ISlotRepository, db *sql.DB) IBookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		slotRepo:    slotRepo,
		db:          db,
	}
}

func (s *BookingService) CreateBooking(userID, slotID uuid.UUID) (*Models.Booking, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	var slot Models.Slot
	querySlot := `SELECT id, room_id, start_time, end_time FROM slots WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(querySlot, slotID).Scan(&slot.ID, &slot.RoomID, &slot.Start, &slot.End)
	if err == sql.ErrNoRows {
		return nil, ErrSlotNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to lock slot: %w", err)
	}

	if slot.Start.Before(time.Now().UTC()) {
		return nil, ErrSlotInPast
	}

	var activeBooking Models.Booking
	queryActive := `SELECT id FROM bookings WHERE slot_id = $1 AND status = 'active' FOR UPDATE`
	err = tx.QueryRow(queryActive, slotID).Scan(&activeBooking.ID)
	if err == nil {
		return nil, ErrSlotAlreadyBooked
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check active booking: %w", err)
	}

	booking := &Models.Booking{
		ID:     uuid.New(),
		SlotID: slotID,
		UserID: userID,
		Status: string(Models.BookingStatusActive),
	}
	insertQuery := `INSERT INTO bookings (id, slot_id, user_id, status) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(insertQuery, booking.ID, booking.SlotID, booking.UserID, booking.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to insert booking: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return booking, nil
}

func (s *BookingService) CancelBooking(bookingID, userID uuid.UUID) (*Models.Booking, error) {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find booking: %w", err)
	}
	if booking == nil {
		return nil, ErrBookingNotFound
	}
	if booking.UserID != userID {
		return nil, ErrNotOwner
	}
	if booking.Status == string(Models.BookingStatusCancelled) {
		return booking, nil
	}

	if err := s.bookingRepo.Cancel(bookingID); err != nil {
		return nil, fmt.Errorf("failed to cancel booking: %w", err)
	}
	booking.Status = string(Models.BookingStatusCancelled)
	return booking, nil
}

func (s *BookingService) GetMyBookings(userID uuid.UUID) ([]Models.Booking, error) {
	bookings, err := s.bookingRepo.FindByUserID(userID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}
	return bookings, nil
}

func (s *BookingService) GetAllBookings(page, pageSize int) ([]Models.Booking, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	bookings, total, err := s.bookingRepo.ListAll(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list bookings: %w", err)
	}
	return bookings, total, nil
}
