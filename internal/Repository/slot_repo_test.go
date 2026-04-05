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

func TestSlotRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewSlotRepository(db)
	slot := &Models.Slot{
		ID:     uuid.New(),
		RoomID: uuid.New(),
		Start:  time.Now().UTC(),
		End:    time.Now().UTC().Add(30 * time.Minute),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO slots (id, room_id, start_time, end_time) VALUES ($1, $2, $3, $4)`)).
		WithArgs(slot.ID, slot.RoomID, slot.Start, slot.End).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(slot)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSlotRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewSlotRepository(db)
	id := uuid.New()
	start := time.Now().UTC()
	end := start.Add(30 * time.Minute)
	rows := sqlmock.NewRows([]string{"id", "room_id", "start_time", "end_time"}).
		AddRow(id, uuid.New(), start, end)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, room_id, start_time, end_time FROM slots WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(rows)

	slot, err := repo.FindByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, slot)
	assert.Equal(t, id, slot.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSlotRepository_DeleteByRoomID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewSlotRepository(db)
	roomID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM slots WHERE room_id = $1`)).
		WithArgs(roomID).
		WillReturnResult(sqlmock.NewResult(0, 3))

	err = repo.DeleteByRoomID(roomID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSlotRepository_FindAvailableByRoomAndDate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewSlotRepository(db)
	roomID := uuid.New()
	date := time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)
	slotID := uuid.New()
	start := time.Date(2026, 4, 5, 9, 0, 0, 0, time.UTC)
	end := start.Add(30 * time.Minute)

	rows := sqlmock.NewRows([]string{"id", "room_id", "start_time", "end_time"}).
		AddRow(slotID, roomID, start, end)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT s.id, s.room_id, s.start_time, s.end_time
		FROM slots s
		LEFT JOIN bookings b ON s.id = b.slot_id AND b.status = 'active'
		WHERE s.room_id = $1
		  AND s.start_time >= $2 AND s.start_time < $3
		  AND b.id IS NULL
		  AND s.start_time >= $4
		ORDER BY s.start_time`)).
		WithArgs(roomID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	slots, err := repo.FindAvailableByRoomAndDate(roomID, date)
	assert.NoError(t, err)
	assert.Len(t, slots, 1)
	assert.Equal(t, slotID, slots[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
