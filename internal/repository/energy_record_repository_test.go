package repository

import (
	"context"
	"database/sql"
	"daya-listrik-api/internal/models"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *EnergyRecordRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	repo := &EnergyRecordRepository{DB: db}
	return db, mock, repo
}

func TestAddRecordSuccess(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	record := &models.EnergyRecord{Usage: 10, Device: "AC", Duration: 2}

	mock.ExpectQuery(`INSERT INTO energy_records`).
		WithArgs(record.Usage, record.Device, record.Duration).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date"}).AddRow(1, time.Now()))

	err := repo.AddRecord(context.Background(), record)

	assert.NoError(t, err)
	assert.Equal(t, 1, record.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIdRecordNotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery(`SELECT id, date, usage, device, duration FROM energy_records WHERE id = \$1`).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	record, err := repo.GetByIdRecord(context.Background(), 99)

	assert.Error(t, err)
	assert.Nil(t, record)
	assert.ErrorIs(t, err, ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteRecordSuccess(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec(`DELETE FROM energy_records WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err := repo.DeleteRecord(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteRecordNotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec(`DELETE FROM energy_records WHERE id = \$1`).
		WithArgs(2).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 row affected

	err := repo.DeleteRecord(context.Background(), 2)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRecordSuccess(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	record := &models.EnergyRecord{ID: 1, Usage: 20, Device: "TV", Duration: 3}

	mock.ExpectQuery(`UPDATE energy_records SET usage=\$1, device=\$2, duration=\$3 WHERE id=\$4 RETURNING date`).
		WithArgs(record.Usage, record.Device, record.Duration, record.ID).
		WillReturnRows(sqlmock.NewRows([]string{"date"}).AddRow(time.Now()))

	err := repo.UpdateRecord(context.Background(), record)

	assert.NoError(t, err)
	assert.False(t, record.Date.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRecordNotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	record := &models.EnergyRecord{ID: 2, Usage: 20, Device: "Lamp", Duration: 1}

	mock.ExpectQuery(`UPDATE energy_records SET usage=\$1, device=\$2, duration=\$3 WHERE id=\$4 RETURNING date`).
		WithArgs(record.Usage, record.Device, record.Duration, record.ID).
		WillReturnError(sql.ErrNoRows)

	err := repo.UpdateRecord(context.Background(), record)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRecordsSuccess(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "date", "usage", "device", "duration"}).
		AddRow(1, time.Now(), 10, "AC", 2).
		AddRow(2, time.Now(), 15, "TV", 3)

	mock.ExpectQuery(`SELECT id, date, usage, device, duration FROM energy_records`).
		WillReturnRows(rows)

	records, err := repo.GetRecords(context.Background())

	assert.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, "AC", records[0].Device)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryDatabaseErrors(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()
	dbErr := errors.New("connection lost")
	mock.ExpectQuery(`INSERT INTO energy_records`).WillReturnError(dbErr)
	assert.ErrorIs(t, repo.AddRecord(context.Background(), &models.EnergyRecord{}), dbErr)
	mock.ExpectQuery(`SELECT id, date, usage, device, duration FROM energy_records WHERE id = \$1`).WithArgs(1).WillReturnError(dbErr)
	_, err := repo.GetByIdRecord(context.Background(), 1)
	assert.ErrorIs(t, err, dbErr)
	mock.ExpectQuery(`SELECT id, date, usage, device, duration FROM energy_records`).WillReturnError(dbErr)
	_, err = repo.GetRecords(context.Background())
	assert.ErrorIs(t, err, dbErr)
	mock.ExpectQuery(`UPDATE energy_records`).WillReturnError(dbErr)
	assert.ErrorIs(t, repo.UpdateRecord(context.Background(), &models.EnergyRecord{ID: 1}), dbErr)
	mock.ExpectExec(`DELETE FROM energy_records`).WillReturnError(dbErr)
	assert.ErrorIs(t, repo.DeleteRecord(context.Background(), 1), dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRecordsEmpty(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()
	mock.ExpectQuery(`SELECT id, date, usage, device, duration FROM energy_records`).WillReturnRows(sqlmock.NewRows([]string{"id", "date", "usage", "device", "duration"}))
	records, err := repo.GetRecords(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, records)
	assert.Empty(t, records)
}
