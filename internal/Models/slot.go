package Models

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	ID     uuid.UUID `db:"id" json:"id"`
	RoomID uuid.UUID `db:"room_id" json:"roomId"`
	Start  time.Time `db:"start_time" json:"start"`
	End    time.Time `db:"end_time" json:"end"`
}
