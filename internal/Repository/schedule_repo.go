package Repository

import (
	"database/sql"
	"fmt"

	"task/internal/Models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type IScheduleRepository interface {
	Create(schedule *Models.Schedule) error
	FindByRoomID(roomID uuid.UUID) (*Models.Schedule, error)
	ExistsForRoom(roomID uuid.UUID) (bool, error)
}

type ScheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) IScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(schedule *Models.Schedule) error {
	query := `INSERT INTO schedules (id, room_id, days_of_week, start_time, end_time) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, schedule.ID, schedule.RoomID, pq.Array(schedule.DaysOfWeek), schedule.StartTime, schedule.EndTime)
	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}
	return nil
}

func (r *ScheduleRepository) FindByRoomID(roomID uuid.UUID) (*Models.Schedule, error) {
	var schedule Models.Schedule
	query := `SELECT id, room_id, days_of_week, start_time, end_time, created_at FROM schedules WHERE room_id = $1`
	var daysOfWeek []int
	err := r.db.QueryRow(query, roomID).Scan(&schedule.ID, &schedule.RoomID, pq.Array(&daysOfWeek), &schedule.StartTime, &schedule.EndTime, &schedule.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find schedule: %w", err)
	}
	schedule.DaysOfWeek = daysOfWeek
	return &schedule, nil
}

func (r *ScheduleRepository) ExistsForRoom(roomID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM schedules WHERE room_id = $1)`
	err := r.db.QueryRow(query, roomID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check schedule existence: %w", err)
	}
	return exists, nil
}
