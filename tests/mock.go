package test

import (
	"context"
	"daya-listrik-api/internal/models"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddRecord(ctx context.Context, record *models.EnergyRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockRepository) GetRecords(ctx context.Context) ([]models.EnergyRecord, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.EnergyRecord), args.Error(1)
}

func (m *MockRepository) DeleteRecord(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) UpdateRecord(ctx context.Context, record *models.EnergyRecord) error {
	if record.ID == 0 {
		return fmt.Errorf("record ID is required for update")
	}
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockRepository) GetByIdRecord(ctx context.Context, id int) (*models.EnergyRecord, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.EnergyRecord), args.Error(1)
}
