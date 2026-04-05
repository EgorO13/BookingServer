package Models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID `db:"id" json:"id"`
	SlotID    uuid.UUID `db:"slot_id" json:"slotId"`
	UserID    uuid.UUID `db:"user_id" json:"userId"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"createdAt,omitempty"`
}

type BookingStatus string

const (
	BookingStatusActive    BookingStatus = "active"
	BookingStatusCancelled BookingStatus = "cancelled"
)
