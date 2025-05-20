package mocks

import (
	"daya-listrik-api/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockEnergyRecordRepository struct {
	mock.Mock // *** embed testify.Mock supaya bisa pakai On, AssertExpectations, dll ***
}

func (m *MockEnergyRecordRepository) AddRecord(record *models.EnergyRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockEnergyRecordRepository) GetRecords() ([]models.EnergyRecord, error) {
	args := m.Called()
	return args.Get(0).([]models.EnergyRecord), args.Error(1)
}

func (m *MockEnergyRecordRepository) GetByIdRecord(id string) (*models.EnergyRecord, error) {
	args := m.Called(id)
	return args.Get(0).(*models.EnergyRecord), args.Error(1)
}

func (m *MockEnergyRecordRepository) UpdateRecord(record *models.EnergyRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockEnergyRecordRepository) DeleteRecord(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
