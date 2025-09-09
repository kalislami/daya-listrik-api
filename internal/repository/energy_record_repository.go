package repository

import (
	"context"
	"database/sql"
	"daya-listrik-api/internal/models"
	"errors"
	"fmt"
)

type EnergyRecordRepositoryInterface interface {
	AddRecord(ctx context.Context, record *models.EnergyRecord) error
	GetByIdRecord(ctx context.Context, id string) (*models.EnergyRecord, error)
	DeleteRecord(ctx context.Context, id string) error
	UpdateRecord(ctx context.Context, record *models.EnergyRecord) error
	GetRecords(ctx context.Context) ([]models.EnergyRecord, error)
}

type EnergyRecordRepository struct {
	DB *sql.DB
}

func (r *EnergyRecordRepository) AddRecord(ctx context.Context, record *models.EnergyRecord) error {
	query := `INSERT INTO energy_records (usage, device, duration) VALUES ($1, $2, $3) RETURNING id, date`
	err := r.DB.QueryRowContext(ctx, query, record.Usage, record.Device, record.Duration).
		Scan(&record.ID, &record.Date)
	if err != nil {
		return fmt.Errorf("error inserting record: %w", err)
	}
	return nil
}

func (r *EnergyRecordRepository) GetByIdRecord(ctx context.Context, id string) (*models.EnergyRecord, error) {
	record := &models.EnergyRecord{}
	query := `SELECT id, date, usage, device, duration FROM energy_records WHERE id = $1`

	err := r.DB.QueryRowContext(ctx, query, id).
		Scan(&record.ID, &record.Date, &record.Usage, &record.Device, &record.Duration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("error retrieving record: %w", err)
	}
	return record, nil
}

func (r *EnergyRecordRepository) DeleteRecord(ctx context.Context, id string) error {
	query := `DELETE FROM energy_records WHERE id = $1`
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record with ID %s not found", id)
	}

	return nil
}

func (r *EnergyRecordRepository) UpdateRecord(ctx context.Context, record *models.EnergyRecord) error {
	query := `UPDATE energy_records SET usage=$1, device=$2, duration=$3 WHERE id=$4`
	result, err := r.DB.ExecContext(ctx, query,
		record.Usage, record.Device, record.Duration, record.ID)
	if err != nil {
		return fmt.Errorf("error updating record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record with ID %d not found", record.ID)
	}

	return nil
}

func (r *EnergyRecordRepository) GetRecords(ctx context.Context) ([]models.EnergyRecord, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, date, usage, device, duration FROM energy_records")
	if err != nil {
		return nil, fmt.Errorf("error fetching records: %w", err)
	}
	defer rows.Close()

	var records []models.EnergyRecord
	for rows.Next() {
		var record models.EnergyRecord
		if err := rows.Scan(&record.ID, &record.Date, &record.Usage, &record.Device, &record.Duration); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in row iteration: %w", err)
	}
	if records == nil {
		records = []models.EnergyRecord{}
	}
	return records, nil
}
