package repository

import (
	"context"
	"daya-listrik-api/internal/models"
	"daya-listrik-api/internal/utils"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EnergyRecordRepositoryInterface interface {
	AddRecord(ctx context.Context, record *models.EnergyRecord) error
	GetByIdRecord(ctx context.Context, id string) (*models.EnergyRecord, error)
	DeleteRecord(ctx context.Context, id string) error
	UpdateRecord(ctx context.Context, record *models.EnergyRecord) error
	GetRecords(ctx context.Context) ([]models.EnergyRecord, error)
}
type EnergyRecordRepository struct {
	DB *gorm.DB
}

func (r *EnergyRecordRepository) AddRecord(ctx context.Context, record *models.EnergyRecord) error {
	// Pakai WithContext biar tetap support context (timeout/cancel)
	if err := r.DB.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("error inserting record: %w", err)
	}
	return nil
}

func (r *EnergyRecordRepository) AddRecordRawQuery(ctx context.Context, record *models.EnergyRecord) error {
	query := `INSERT INTO energy_records (usage, device, duration) 
	          VALUES (?, ?, ?) RETURNING id, date`

	// Raw query dengan Scan ke struct
	if err := r.DB.WithContext(ctx).
		Raw(query, record.Usage, record.Device, record.Duration).
		Scan(&record).Error; err != nil {
		return fmt.Errorf("error inserting record: %w", err)
	}
	return nil
}

func (r *EnergyRecordRepository) GetByIdRecord(ctx context.Context, id string) (*models.EnergyRecord, error) {
	record := &models.EnergyRecord{}

	// ORM style pakai Where + First
	if err := r.DB.WithContext(ctx).
		Where("id = ?", id).
		First(record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(fiber.StatusNotFound, "record not found")
		}
		return nil, fmt.Errorf("error retrieving record: %w", err)
	}

	return record, nil
}

func (r *EnergyRecordRepository) GetByIdRecordRawQuery(ctx context.Context, id string) (*models.EnergyRecord, error) {
	record := &models.EnergyRecord{}
	query := `SELECT id, date, usage, device, duration FROM energy_records WHERE id = ?`

	if err := r.DB.WithContext(ctx).
		Raw(query, id).
		Scan(record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("error retrieving record: %w", err)
	}

	return record, nil
}

func (r *EnergyRecordRepository) DeleteRecord(ctx context.Context, id string) error {
	result := r.DB.WithContext(ctx).Delete(&models.EnergyRecord{}, id)

	if result.Error != nil {
		return fmt.Errorf("error deleting record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("record with ID %s not found", id)
	}

	return nil
}

func (r *EnergyRecordRepository) DeleteRecordRawQuery(ctx context.Context, id string) error {
	// pakai ? sebagai placeholder biar portable ke semua DB
	result := r.DB.WithContext(ctx).Exec("DELETE FROM energy_records WHERE id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error deleting record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("record with ID %s not found", id)
	}

	return nil
}

func (r *EnergyRecordRepository) UpdateRecord(ctx context.Context, record *models.EnergyRecord) error {
	result := r.DB.WithContext(ctx).
		Model(&models.EnergyRecord{}).
		Where("id = ?", record.ID).
		Updates(map[string]interface{}{
			"usage":    record.Usage,
			"device":   record.Device,
			"duration": record.Duration,
		})

	if result.Error != nil {
		return fmt.Errorf("error updating record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return utils.NewAppError(fiber.StatusNotFound, "record not found")
	}

	return nil
}

func (r *EnergyRecordRepository) UpdateRecordRawQuery(ctx context.Context, record *models.EnergyRecord) error {
	result := r.DB.WithContext(ctx).Exec(
		"UPDATE energy_records SET usage = ?, device = ?, duration = ? WHERE id = ?",
		record.Usage, record.Device, record.Duration, record.ID,
	)

	if result.Error != nil {
		return fmt.Errorf("error updating record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("record with ID %d not found", record.ID)
	}

	return nil
}

func (r *EnergyRecordRepository) GetRecords(ctx context.Context) ([]models.EnergyRecord, error) {
	var records []models.EnergyRecord

	result := r.DB.WithContext(ctx).Find(&records)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching records: %w", result.Error)
	}

	// biar aman, jangan return nil slice
	if records == nil {
		records = []models.EnergyRecord{}
	}

	return records, nil
}

func (r *EnergyRecordRepository) GetRecordsRawQuery(ctx context.Context) ([]models.EnergyRecord, error) {
	var records []models.EnergyRecord

	rows, err := r.DB.WithContext(ctx).Raw(
		"SELECT id, date, usage, device, duration FROM energy_records",
	).Rows()
	if err != nil {
		return nil, fmt.Errorf("error fetching records: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record models.EnergyRecord
		if err := rows.Scan(&record.ID, &record.Date, &record.Usage, &record.Device, &record.Duration); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		records = append(records, record)
	}

	if records == nil {
		records = []models.EnergyRecord{}
	}

	return records, nil
}
