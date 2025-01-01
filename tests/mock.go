package test

import (
	"bytes"
	"daya-listrik-api/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddRecord(record *models.EnergyRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockRepository) GetRecords() ([]models.EnergyRecord, error) {
	args := m.Called()
	return args.Get(0).([]models.EnergyRecord), args.Error(1)
}

func (m *MockRepository) DeleteRecord(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) UpdateRecord(record *models.EnergyRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockRepository) GetByIdRecord(id string) (*models.EnergyRecord, error) {
	args := m.Called(id)
	return args.Get(0).(*models.EnergyRecord), args.Error(1)
}

func MakeRequest(method, url string, body []byte) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	return req, rr
}

func AssertStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	if status := rr.Code; status != expected {
		t.Errorf("Expected status code %v, got %v", expected, status)
	}
}