package Repository

import (
	"database/sql"
	"fmt"

	"task/internal/Models"

	"github.com/google/uuid"
)

type IRoomRepository interface {
	Create(room *Models.Room) error
	List() ([]Models.Room, error)
	Exists(id uuid.UUID) (bool, error)
}

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) IRoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *Models.Room) error {
	query := `INSERT INTO rooms (id, name, description, capacity) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, room.ID, room.Name, room.Description, room.Capacity)
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}
	return nil
}

func (r *RoomRepository) List() ([]Models.Room, error) {
	rows, err := r.db.Query(`SELECT id, name, description, capacity, created_at FROM rooms ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list rooms: %w", err)
	}
	defer rows.Close()

	var rooms []Models.Room
	for rows.Next() {
		var room Models.Room
		err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan room: %w", err)
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *RoomRepository) Exists(id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)`
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check room existence: %w", err)
	}
	return exists, nil
}
