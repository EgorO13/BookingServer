package Models

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID         uuid.UUID `db:"id" json:"id"`
	RoomID     uuid.UUID `db:"room_id" json:"roomId"`
	DaysOfWeek []int     `db:"days_of_week" json:"daysOfWeek"`
	StartTime  string    `db:"start_time" json:"startTime"`
	EndTime    string    `db:"end_time" json:"endTime"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt,omitempty"`
}
