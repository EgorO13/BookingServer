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

func TestBookingRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewBookingRepository(db)
	booking := &Models.Booking{
		ID:     uuid.New(),
		SlotID: uuid.New(),
		UserID: uuid.New(),
		Status: string(Models.BookingStatusActive),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO bookings (id, slot_id, user_id, status) VALUES ($1, $2, $3, $4)`)).
		WithArgs(booking.ID, booking.SlotID, booking.UserID, booking.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(booking)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookingRepository_Cancel(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewBookingRepository(db)
	id := uuid.New()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE bookings SET status = 'cancelled' WHERE id = $1 AND status = 'active'`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	err = repo.Cancel(id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookingRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewBookingRepository(db)
	id := uuid.New()
	createdAt := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "slot_id", "user_id", "status", "created_at"}).
		AddRow(id, uuid.New(), uuid.New(), "active", createdAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, slot_id, user_id, status, created_at FROM bookings WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(rows)

	booking, err := repo.FindByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.Equal(t, id, booking.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookingRepository_FindActiveBySlotID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewBookingRepository(db)
	slotID := uuid.New()
	createdAt := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "slot_id", "user_id", "status", "created_at"}).
		AddRow(uuid.New(), slotID, uuid.New(), "active", createdAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, slot_id, user_id, status, created_at FROM bookings WHERE slot_id = $1 AND status = 'active'`)).
		WithArgs(slotID).
		WillReturnRows(rows)

	booking, err := repo.FindActiveBySlotID(slotID)
	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookingRepository_FindByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewBookingRepository(db)
	userID := uuid.New()
	createdAt := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "slot_id", "user_id", "status", "created_at"}).
		AddRow(uuid.New(), uuid.New(), userID, "active", createdAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.id, b.slot_id, b.user_id, b.status, b.created_at FROM bookings b JOIN slots s ON b.slot_id = s.id WHERE b.user_id = $1 ORDER BY s.start_time ASC`)).
		WithArgs(userID).
		WillReturnRows(rows)

	bookings, err := repo.FindByUserID(userID, false)
	assert.NoError(t, err)
	assert.Len(t, bookings, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookingRepository_ListAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewBookingRepository(db)
	createdAt := time.Now().UTC()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM bookings`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	rows := sqlmock.NewRows([]string{"id", "slot_id", "user_id", "status", "created_at"}).
		AddRow(uuid.New(), uuid.New(), uuid.New(), "active", createdAt).
		AddRow(uuid.New(), uuid.New(), uuid.New(), "cancelled", createdAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, slot_id, user_id, status, created_at FROM bookings ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
		WithArgs(20, 0).
		WillReturnRows(rows)

	bookings, total, err := repo.ListAll(1, 20)
	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	assert.Equal(t, 2, total)
	assert.NoError(t, mock.ExpectationsWereMet())
}
