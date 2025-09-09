package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"daya-listrik-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEnergyRecordRepo struct {
	mock.Mock
}

const (
	pathWithId            = "/api/records/1"
	HeaderContentType     = "Content-Type"
	HeaderApplicationJSON = "application/json"
)

func (m *MockEnergyRecordRepo) AddRecord(ctx context.Context, record *models.EnergyRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockEnergyRecordRepo) GetRecords(ctx context.Context) ([]models.EnergyRecord, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.EnergyRecord), args.Error(1)
}

func (m *MockEnergyRecordRepo) DeleteRecord(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEnergyRecordRepo) UpdateRecord(ctx context.Context, record *models.EnergyRecord) error {
	if record.ID == 0 {
		return fmt.Errorf("record ID is required for update")
	}
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockEnergyRecordRepo) GetByIdRecord(ctx context.Context, id string) (*models.EnergyRecord, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.EnergyRecord), args.Error(1)
}

// Helper untuk setup Fiber app + handler
func setupTestApp() (*fiber.App, *MockEnergyRecordRepo) {
	app := fiber.New()
	mockRepo := new(MockEnergyRecordRepo)
	handler := &EnergyRecordHandler{Repo: mockRepo}
	InitializeRoutes(app, handler)
	return app, mockRepo
}

func TestAddRecordSuccess(t *testing.T) {
	app, mockRepo := setupTestApp()

	payload := `{"usage": 10, "device": "AC"}`
	record := &models.EnergyRecord{Usage: 10, Device: "AC"}

	mockRepo.On("AddRecord", mock.Anything, record).Return(nil)

	req := httptest.NewRequest("POST", "/api/records/add", strings.NewReader(payload))
	req.Header.Set(HeaderContentType, HeaderApplicationJSON)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var resRecord models.EnergyRecord
	json.NewDecoder(resp.Body).Decode(&resRecord)
	assert.Equal(t, record.Usage, resRecord.Usage)
	assert.Equal(t, record.Device, resRecord.Device)
}

func TestAddRecordInvalidBody(t *testing.T) {
	app, _ := setupTestApp()

	payload := `{"usage": -5, "device": ""}` // invalid

	req := httptest.NewRequest("POST", "/api/records/add", strings.NewReader(payload))
	req.Header.Set(HeaderContentType, HeaderApplicationJSON)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetRecordsSuccess(t *testing.T) {
	app, mockRepo := setupTestApp()

	mockData := []models.EnergyRecord{
		{ID: 1, Usage: 10, Device: "AC"},
		{ID: 2, Usage: 5, Device: "Fan"},
	}
	mockRepo.On("GetRecords", mock.Anything).Return(mockData, nil)

	req := httptest.NewRequest("GET", "/api/records", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var records []models.EnergyRecord
	json.NewDecoder(resp.Body).Decode(&records)
	assert.Len(t, records, 2)
	assert.Equal(t, "AC", records[0].Device)
}

func TestDeleteRecordInvalidId(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("DELETE", "/api/records/abc", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteRecordSuccess(t *testing.T) {
	app, mockRepo := setupTestApp()

	mockRepo.On("DeleteRecord", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest("DELETE", pathWithId, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestGetByIdRecordNotFound(t *testing.T) {
	app, mockRepo := setupTestApp()

	mockRepo.On("GetByIdRecord", mock.Anything, "1").Return((*models.EnergyRecord)(nil), errors.New("not found"))

	req := httptest.NewRequest("GET", pathWithId, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetByIdRecordSuccess(t *testing.T) {
	app, mockRepo := setupTestApp()

	mockRecord := &models.EnergyRecord{ID: 1, Usage: 15, Device: "Heater"}
	mockRepo.On("GetByIdRecord", mock.Anything, "1").Return(mockRecord, nil)

	req := httptest.NewRequest("GET", pathWithId, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var record models.EnergyRecord
	json.NewDecoder(resp.Body).Decode(&record)
	assert.Equal(t, "Heater", record.Device)
}

func TestUpdateRecordSuccess(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockEnergyRecordRepo)
	handler := &EnergyRecordHandler{Repo: mockRepo}
	InitializeRoutes(app, handler)

	payload := `{"usage": 10, "device": "AC"}`
	record := &models.EnergyRecord{ID: 1, Usage: 10, Device: "AC"}

	mockRepo.On("UpdateRecord", mock.Anything, record).Return(nil)

	req := httptest.NewRequest("PUT", pathWithId, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var resRecord models.EnergyRecord
	json.NewDecoder(resp.Body).Decode(&resRecord)
	assert.Equal(t, record.ID, resRecord.ID)
	assert.Equal(t, record.Device, resRecord.Device)
}

func TestUpdateRecordInvalidParamId(t *testing.T) {
	app, _ := setupTestApp()

	payload := `{"usage": 10, "device": "AC"}`
	req := httptest.NewRequest("PUT", "/api/records/abc", strings.NewReader(payload))
	req.Header.Set(HeaderContentType, HeaderApplicationJSON)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateRecordInvalidBody(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockEnergyRecordRepo)
	handler := &EnergyRecordHandler{Repo: mockRepo}
	InitializeRoutes(app, handler)

	payload := `{"usage": -5, "device": ""}` // invalid usage/device
	req := httptest.NewRequest("PUT", pathWithId, strings.NewReader(payload))
	req.Header.Set(HeaderContentType, HeaderApplicationJSON)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateRecordRepoError(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockEnergyRecordRepo)
	handler := &EnergyRecordHandler{Repo: mockRepo}
	InitializeRoutes(app, handler)

	payload := `{"usage": 10, "device": "AC"}`
	record := &models.EnergyRecord{ID: 1, Usage: 10, Device: "AC"}

	mockRepo.On("UpdateRecord", mock.Anything, record).Return(fmt.Errorf("repo error"))

	req := httptest.NewRequest("PUT", pathWithId, strings.NewReader(payload))
	req.Header.Set(HeaderContentType, HeaderApplicationJSON)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
