package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"daya-listrik-api/internal/models"
	"daya-listrik-api/internal/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *EnergyRecordRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := &EnergyRecordRepository{DB: gormDB}
	return gormDB, mock, repo
}

func TestAddRecordSuccess(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	record := &models.EnergyRecord{Usage: 10, Device: "AC", Duration: 2}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `energy_records` (`date`,`usage`,`device`,`duration`) VALUES (?,?,?,?)")).
		WithArgs(sqlmock.AnyArg(), record.Usage, record.Device, record.Duration).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.AddRecord(context.Background(), record)

	assert.NoError(t, err)
	assert.Equal(t, int(1), record.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIdRecordNotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `energy_records` WHERE id = ? ORDER BY `energy_records`.`id` LIMIT ?")).
		WithArgs("99", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.GetByIdRecord(context.Background(), "99")

	appErr, ok := err.(*utils.AppError)
	assert.True(t, ok, "expected error to be of type *utils.AppError")
	assert.Equal(t, 404, appErr.Code)
	assert.Equal(t, "record not found", appErr.Message)
}

func TestDeleteRecordSuccess(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `energy_records` WHERE `energy_records`.`id` = ?")).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteRecord(context.Background(), "1")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteRecordNotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `energy_records` WHERE `energy_records`.`id` = ?")).
		WithArgs("2").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DeleteRecord(context.Background(), "2")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRecordSuccess(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE `energy_records` SET `device`=?,`duration`=?,`usage`=? WHERE id = ?",
	)).
		WithArgs("TV", 3.0, 20.0, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	record := &models.EnergyRecord{ID: 1, Device: "TV", Duration: 3, Usage: 20}
	err := repo.UpdateRecord(context.Background(), record)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRecordNotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE `energy_records` SET `device`=?,`duration`=?,`usage`=? WHERE id = ?",
	)).
		WithArgs("Lamp", 1.0, 20.0, 2).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
	mock.ExpectCommit()

	record := &models.EnergyRecord{ID: 2, Device: "Lamp", Duration: 1.0, Usage: 20.0}
	err := repo.UpdateRecord(context.Background(), record)

	appErr, ok := err.(*utils.AppError)
	assert.True(t, ok, "expected error to be of type *utils.AppError")
	assert.Equal(t, 404, appErr.Code)
	assert.Equal(t, "record not found", appErr.Message)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRecordsSuccess(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"id", "date", "usage", "device", "duration"}).
		AddRow(1, time.Now(), 10, "AC", 2).
		AddRow(2, time.Now(), 15, "TV", 3)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `energy_records`")).
		WillReturnRows(rows)

	records, err := repo.GetRecords(context.Background())

	assert.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, "AC", records[0].Device)
	assert.NoError(t, mock.ExpectationsWereMet())
}
