package handlers

import (
	"bytes"
	"daya-listrik-api/internal/models"
	"daya-listrik-api/internal/repository/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddRecord_Success(t *testing.T) {
	mockRepo := new(mocks.MockEnergyRecordRepository)
	handler := AddRecord(mockRepo)

	record := models.EnergyRecord{
		Device: "AC",
		Usage:  100,
	}
	body, _ := json.Marshal(record)

	mockRepo.On("AddRecord", mock.AnythingOfType("*models.EnergyRecord")).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/records/add", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetRecords_Success(t *testing.T) {
	mockRepo := new(mocks.MockEnergyRecordRepository)
	handler := GetRecords(mockRepo)

	records := []models.EnergyRecord{
		{ID: 1, Device: "Lamp", Usage: 20},
	}
	mockRepo.On("GetRecords").Return(records, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/records", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []models.EnergyRecord
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, records, resp)
	mockRepo.AssertExpectations(t)
}

func TestGetByIdRecords_Success(t *testing.T) {
	mockRepo := new(mocks.MockEnergyRecordRepository)
	handler := GetByIdRecords(mockRepo)

	record := &models.EnergyRecord{ID: 1, Device: "TV", Usage: 50}
	mockRepo.On("GetByIdRecord", "1").Return(record, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/records/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.EnergyRecord
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, *record, resp)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRecords_Success(t *testing.T) {
	mockRepo := new(mocks.MockEnergyRecordRepository)
	handler := UpdateRecords(mockRepo)

	record := models.EnergyRecord{ID: 1, Device: "Fan", Usage: 60}
	body, _ := json.Marshal(record)

	mockRepo.On("UpdateRecord", mock.AnythingOfType("*models.EnergyRecord")).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/api/records/1", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.EnergyRecord
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, record, resp)
	mockRepo.AssertExpectations(t)
}

func TestDeleteRecords_Success(t *testing.T) {
	mockRepo := new(mocks.MockEnergyRecordRepository)
	handler := DeleteRecords(mockRepo)

	mockRepo.On("DeleteRecord", "1").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/records/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockRepo.AssertExpectations(t)
}
