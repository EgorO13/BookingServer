package Repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewUserRepository(db)
	userID := uuid.New()
	email := "test@example.com"
	role := "user"
	createdAt := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
		AddRow(userID, email, role, createdAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, role, created_at FROM users WHERE id = $1`)).
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.FindByID(userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, role, user.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewUserRepository(db)
	userID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, role, created_at FROM users WHERE id = $1`)).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindByID(userID)
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}
