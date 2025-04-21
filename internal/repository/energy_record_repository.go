package repository

import (
	"database/sql"
	"daya-listrik-api/internal/models"
	"fmt"
)

type EnergyRecordRepositoryInterface interface {
	AddRecord(record *models.EnergyRecord) error
	GetByIdRecord(id string) (*models.EnergyRecord, error)
	DeleteRecord(id string) error
	UpdateRecord(record *models.EnergyRecord) error
	GetRecords() ([]models.EnergyRecord, error)
}

type EnergyRecordRepository struct {
	DB *sql.DB
}

func (r *EnergyRecordRepository) AddRecord(record *models.EnergyRecord) error {
	query := `INSERT INTO energy_records (usage, device, duration) VALUES ($1, $2, $3) RETURNING id, date`
	err := r.DB.QueryRow(query, record.Usage, record.Device, record.Duration).Scan(&record.ID, &record.Date)
	if err != nil {
		return fmt.Errorf("error inserting record: %v", err)
	}
	return nil
}

func (r *EnergyRecordRepository) GetByIdRecord(id string) (*models.EnergyRecord, error) {
	record := &models.EnergyRecord{}
	query := `SELECT id, date, usage, device, duration FROM energy_records WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&record.ID, &record.Date, &record.Usage, &record.Device, &record.Duration)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.EnergyRecord{}, fmt.Errorf("record with ID %s not found", id)
		}
		return &models.EnergyRecord{}, fmt.Errorf("error retrieving record: %v", err)
	}
	return record, nil
}

func (r *EnergyRecordRepository) DeleteRecord(id string) error {
	query := `DELETE FROM energy_records WHERE id = $1`
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting record: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record with ID %s not found", id)
	}

	return nil
}

func (r *EnergyRecordRepository) UpdateRecord(record *models.EnergyRecord) error {
	query := `UPDATE energy_records SET usage=$1, device=$2, duration=$3 WHERE id=$4`
	result, err := r.DB.Exec(query, record.Usage, record.Device, record.Duration, record.ID)
	if err != nil {
		return fmt.Errorf("error updating record: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record with ID %d not found", record.ID)
	}

	return nil
}

func (r *EnergyRecordRepository) GetRecords() ([]models.EnergyRecord, error) {
	rows, err := r.DB.Query("SELECT id, date, usage, device, duration FROM energy_records")
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
