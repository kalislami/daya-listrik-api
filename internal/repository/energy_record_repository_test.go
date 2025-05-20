package repository

import (
	"database/sql"
	"daya-listrik-api/internal/models"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &EnergyRecordRepository{DB: db}

	record := &models.EnergyRecord{
		Usage:    10.5,
		Device:   "Device A",
		Duration: 5.0,
	}

	// Setup expectation for QueryRow + Scan on Insert returning id and date
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO energy_records (usage, device, duration) VALUES ($1, $2, $3) RETURNING id, date`,
	)).WithArgs(record.Usage, record.Device, record.Duration).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date"}).AddRow(1, time.Now()))

	err = repo.AddRecord(record)
	assert.NoError(t, err)
	assert.Equal(t, 1, record.ID)

	// Test error on Insert
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO energy_records (usage, device, duration) VALUES ($1, $2, $3) RETURNING id, date`,
	)).WithArgs(record.Usage, record.Device, record.Duration).
		WillReturnError(errors.New("insert error"))

	err = repo.AddRecord(record)
	assert.Error(t, err)
}

func TestGetByIdRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &EnergyRecordRepository{DB: db}

	id := "1"
	expectedRecord := &models.EnergyRecord{
		ID:       1,
		Date:     time.Now(),
		Usage:    20.0,
		Device:   "Device X",
		Duration: 3.5,
	}

	// Happy path
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, date, usage, device, duration FROM energy_records WHERE id = $1`,
	)).WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date", "usage", "device", "duration"}).
			AddRow(expectedRecord.ID, expectedRecord.Date, expectedRecord.Usage, expectedRecord.Device, expectedRecord.Duration))

	rec, err := repo.GetByIdRecord(id)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, rec.ID)

	// No rows found
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, date, usage, device, duration FROM energy_records WHERE id = $1`,
	)).WithArgs("999").
		WillReturnError(sql.ErrNoRows)

	rec, err = repo.GetByIdRecord("999")
	assert.Error(t, err)
	assert.Equal(t, 0, rec.ID)

	// Other error
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, date, usage, device, duration FROM energy_records WHERE id = $1`,
	)).WithArgs("error").
		WillReturnError(errors.New("some db error"))

	rec, err = repo.GetByIdRecord("error")
	assert.Error(t, err)
}

func TestDeleteRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &EnergyRecordRepository{DB: db}

	id := "1"

	// Successful delete (rows affected = 1)
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM energy_records WHERE id = $1`,
	)).WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteRecord(id)
	assert.NoError(t, err)

	// Delete no rows affected
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM energy_records WHERE id = $1`,
	)).WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.DeleteRecord(id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Exec error
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM energy_records WHERE id = $1`,
	)).WithArgs(id).
		WillReturnError(errors.New("exec error"))

	err = repo.DeleteRecord(id)
	assert.Error(t, err)
}

func TestUpdateRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &EnergyRecordRepository{DB: db}

	record := &models.EnergyRecord{
		ID:       1,
		Usage:    15.0,
		Device:   "Device B",
		Duration: 2.5,
	}

	// Successful update (rows affected = 1)
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE energy_records SET usage=$1, device=$2, duration=$3 WHERE id=$4`,
	)).WithArgs(record.Usage, record.Device, record.Duration, record.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateRecord(record)
	assert.NoError(t, err)

	// Update no rows affected
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE energy_records SET usage=$1, device=$2, duration=$3 WHERE id=$4`,
	)).WithArgs(record.Usage, record.Device, record.Duration, record.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateRecord(record)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Exec error
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE energy_records SET usage=$1, device=$2, duration=$3 WHERE id=$4`,
	)).WithArgs(record.Usage, record.Device, record.Duration, record.ID).
		WillReturnError(errors.New("exec error"))

	err = repo.UpdateRecord(record)
	assert.Error(t, err)
}

func TestGetRecords(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &EnergyRecordRepository{DB: db}

	// Happy path with 2 records
	rows := sqlmock.NewRows([]string{"id", "date", "usage", "device", "duration"}).
		AddRow(1, time.Now(), 10.0, "Device1", 2.0).
		AddRow(2, time.Now(), 20.0, "Device2", 3.0)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, date, usage, device, duration FROM energy_records`,
	)).WillReturnRows(rows)

	records, err := repo.GetRecords()
	assert.NoError(t, err)
	assert.Len(t, records, 2)

	// Empty result
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, date, usage, device, duration FROM energy_records`,
	)).WillReturnRows(sqlmock.NewRows([]string{"id", "date", "usage", "device", "duration"}))

	records, err = repo.GetRecords()
	assert.NoError(t, err)
	assert.Len(t, records, 0)

	// Query error
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, date, usage, device, duration FROM energy_records`,
	)).WillReturnError(errors.New("query error"))

	records, err = repo.GetRecords()
	assert.Error(t, err)
}

func TestGetRecordsScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &EnergyRecordRepository{DB: db}

	rows := sqlmock.NewRows([]string{"id", "date", "usage", "device", "duration"}).
		AddRow(1, time.Now(), "invalid_float", "Device1", 2.0)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, date, usage, device, duration FROM energy_records`,
	)).WillReturnRows(rows)

	records, err := repo.GetRecords()
	assert.Error(t, err)
	assert.Nil(t, records)
}

func TestGetRecordsRowsErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &EnergyRecordRepository{DB: db}

	rows := sqlmock.NewRows([]string{"id", "date", "usage", "device", "duration"}).
		AddRow(1, time.Now(), 10.0, "Device1", 2.0)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, date, usage, device, duration FROM energy_records`,
	)).WillReturnRows(rows)

	records, err := repo.GetRecords()
	assert.NoError(t, err)
	assert.Len(t, records, 1)
}