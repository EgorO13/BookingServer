package Repository

import (
	"database/sql"
	"fmt"

	"task/internal/Models"

	"github.com/google/uuid"
)

type IUserRepository interface {
	FindByID(id uuid.UUID) (*Models.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uuid.UUID) (*Models.User, error) {
	var user Models.User
	query := `SELECT id, email, role, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}
