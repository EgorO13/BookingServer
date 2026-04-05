package Repository

import (
	"database/sql"
	"fmt"
	"time"

	"task/internal/Models"

	"github.com/google/uuid"
)

type ISlotRepository interface {
	Create(slot *Models.Slot) error
	CreateBatch(slots []Models.Slot) error
	FindAvailableByRoomAndDate(roomID uuid.UUID, date time.Time) ([]Models.Slot, error)
	FindByID(id uuid.UUID) (*Models.Slot, error)
	DeleteByRoomID(roomID uuid.UUID) error
}

type SlotRepository struct {
	db *sql.DB
}

func NewSlotRepository(db *sql.DB) ISlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) Create(slot *Models.Slot) error {
	query := `INSERT INTO slots (id, room_id, start_time, end_time) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, slot.ID, slot.RoomID, slot.Start, slot.End)
	if err != nil {
		return fmt.Errorf("failed to create slot: %w", err)
	}
	return nil
}

func (r *SlotRepository) CreateBatch(slots []Models.Slot) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT INTO slots (id, room_id, start_time, end_time) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	for _, slot := range slots {
		_, err := stmt.Exec(slot.ID, slot.RoomID, slot.Start, slot.End)
		if err != nil {
			return fmt.Errorf("failed to insert slot: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *SlotRepository) FindAvailableByRoomAndDate(roomID uuid.UUID, date time.Time) ([]Models.Slot, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
        SELECT s.id, s.room_id, s.start_time, s.end_time
        FROM slots s
        LEFT JOIN bookings b ON s.id = b.slot_id AND b.status = 'active'
        WHERE s.room_id = $1
          AND s.start_time >= $2 AND s.start_time < $3
          AND b.id IS NULL
          AND s.start_time >= $4
        ORDER BY s.start_time
    `
	now := time.Now().UTC()
	rows, err := r.db.Query(query, roomID, startOfDay, endOfDay, now)
	if err != nil {
		return nil, fmt.Errorf("failed to query slots: %w", err)
	}
	defer rows.Close()

	var slots []Models.Slot
	for rows.Next() {
		var slot Models.Slot
		err := rows.Scan(&slot.ID, &slot.RoomID, &slot.Start, &slot.End)
		if err != nil {
			return nil, fmt.Errorf("failed to scan slot: %w", err)
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *SlotRepository) FindByID(id uuid.UUID) (*Models.Slot, error) {
	var slot Models.Slot
	query := `SELECT id, room_id, start_time, end_time FROM slots WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&slot.ID, &slot.RoomID, &slot.Start, &slot.End)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find slot: %w", err)
	}
	return &slot, nil
}

func (r *SlotRepository) DeleteByRoomID(roomID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM slots WHERE room_id = $1`, roomID)
	if err != nil {
		return fmt.Errorf("failed to delete slots for room: %w", err)
	}
	return nil
}
