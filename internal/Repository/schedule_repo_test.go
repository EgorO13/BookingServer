package Repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"task/internal/Models"
)

func TestScheduleRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewScheduleRepository(db)
	schedule := &Models.Schedule{
		ID:         uuid.New(),
		RoomID:     uuid.New(),
		DaysOfWeek: []int{1, 2, 3, 4, 5},
		StartTime:  "09:00",
		EndTime:    "17:00",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO schedules (id, room_id, days_of_week, start_time, end_time) VALUES ($1, $2, $3, $4, $5)`)).
		WithArgs(schedule.ID, schedule.RoomID, pq.Array(schedule.DaysOfWeek), schedule.StartTime, schedule.EndTime).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(schedule)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleRepository_ExistsForRoom(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewScheduleRepository(db)
	roomID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM schedules WHERE room_id = $1)`)).
		WithArgs(roomID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsForRoom(roomID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}
