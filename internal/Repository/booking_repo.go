package Repository

import (
	"database/sql"
	"fmt"
	"time"

	"task/internal/Models"

	"github.com/google/uuid"
)

type IBookingRepository interface {
	Create(booking *Models.Booking) error
	Cancel(id uuid.UUID) error
	FindByID(id uuid.UUID) (*Models.Booking, error)
	FindActiveBySlotID(slotID uuid.UUID) (*Models.Booking, error)
	FindByUserID(userID uuid.UUID, onlyFuture bool) ([]Models.Booking, error)
	ListAll(page, pageSize int) ([]Models.Booking, int, error)
}

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) IBookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *Models.Booking) error {
	query := `INSERT INTO bookings (id, slot_id, user_id, status) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, booking.ID, booking.SlotID, booking.UserID, booking.Status)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}
	return nil
}

func (r *BookingRepository) Cancel(id uuid.UUID) error {
	query := `UPDATE bookings SET status = 'cancelled' WHERE id = $1 AND status = 'active'`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}
	return nil
}

func (r *BookingRepository) FindByID(id uuid.UUID) (*Models.Booking, error) {
	var booking Models.Booking
	query := `SELECT id, slot_id, user_id, status, created_at FROM bookings WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&booking.ID, &booking.SlotID, &booking.UserID, &booking.Status, &booking.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find booking: %w", err)
	}
	return &booking, nil
}

func (r *BookingRepository) FindActiveBySlotID(slotID uuid.UUID) (*Models.Booking, error) {
	var booking Models.Booking
	query := `SELECT id, slot_id, user_id, status, created_at FROM bookings WHERE slot_id = $1 AND status = 'active'`
	err := r.db.QueryRow(query, slotID).Scan(&booking.ID, &booking.SlotID, &booking.UserID, &booking.Status, &booking.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find active booking for slot: %w", err)
	}
	return &booking, nil
}

func (r *BookingRepository) FindByUserID(userID uuid.UUID, onlyFuture bool) ([]Models.Booking, error) {
	query := `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.created_at
		FROM bookings b
		JOIN slots s ON b.slot_id = s.id
		WHERE b.user_id = $1
	`
	args := []interface{}{userID}

	if onlyFuture {
		query += ` AND s.start_time >= $2`
		args = append(args, time.Now().UTC())
	}
	query += ` ORDER BY s.start_time ASC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	var bookings []Models.Booking
	for rows.Next() {
		var booking Models.Booking
		err := rows.Scan(&booking.ID, &booking.SlotID, &booking.UserID, &booking.Status, &booking.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (r *BookingRepository) ListAll(page, pageSize int) ([]Models.Booking, int, error) {
	offset := (page - 1) * pageSize

	var total int
	countQuery := `SELECT COUNT(*) FROM bookings`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	query := `
        SELECT id, slot_id, user_id, status, created_at
        FROM bookings
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	var bookings []Models.Booking
	for rows.Next() {
		var booking Models.Booking
		err := rows.Scan(&booking.ID, &booking.SlotID, &booking.UserID, &booking.Status, &booking.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}
	return bookings, total, nil
}
