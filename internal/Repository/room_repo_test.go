package Repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"task/internal/Models"
)

func TestRoomRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewRoomRepository(db)
	room := &Models.Room{
		ID:          uuid.New(),
		Name:        "Test Room",
		Description: stringPtr("desc"),
		Capacity:    intPtr(10),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO rooms (id, name, description, capacity) VALUES ($1, $2, $3, $4)`)).
		WithArgs(room.ID, room.Name, room.Description, room.Capacity).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(room)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewRoomRepository(db)
	id1 := uuid.New()
	id2 := uuid.New()
	now := time.Now().UTC()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "capacity", "created_at"}).
		AddRow(id1, "Room1", nil, nil, now).
		AddRow(id2, "Room2", stringPtr("desc"), intPtr(5), now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, description, capacity, created_at FROM rooms ORDER BY created_at DESC`)).
		WillReturnRows(rows)

	rooms, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, rooms, 2)
	assert.Equal(t, id1, rooms[0].ID)
	assert.Equal(t, "Room1", rooms[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomRepository_Exists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewRoomRepository(db)
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.Exists(id)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func stringPtr(s string) *string { return &s }
func intPtr(i int) *int          { return &i }
