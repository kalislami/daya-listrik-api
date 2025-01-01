package test

import (
	"daya-listrik-api/internal/models"
	"fmt"
	"time"

	"daya-listrik-api/internal/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAddRecord(t *testing.T) {
	const routeApi = "/api/records/add"

	t.Run("should get status 201: success add record", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRecord := &models.EnergyRecord{
			Usage:  100,
			Device: "Laptop",
		}

		mockRepo.On("AddRecord", mockRecord).Return(nil)

		body, _ := json.Marshal(mockRecord)
		req, rr := MakeRequest("POST", routeApi, body)

		handler := handlers.AddRecord(mockRepo)
		handler.ServeHTTP(rr, req)

		AssertStatusCode(t, rr, http.StatusCreated)

		var response models.EnergyRecord
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}
		assert.Equal(t, mockRecord, &response)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should get status 400: invalid request", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRecord := &models.EnergyRecord{Usage: 100}

		body, _ := json.Marshal(mockRecord)
		req, rr := MakeRequest("POST", routeApi, body)

		handler := handlers.AddRecord(mockRepo)
		handler.ServeHTTP(rr, req)

		AssertStatusCode(t, rr, http.StatusBadRequest)

		mockRepo.AssertNotCalled(t, "AddRecord")
	})

	t.Run("should return status 400: invalid JSON body", func(t *testing.T) {
		mockRepo := new(MockRepository)
		invalidPayload := `{"Usage": "invalid_number", "Device": 123}`

		req, rr := MakeRequest("POST", routeApi, []byte(invalidPayload))

		handler := handlers.AddRecord(mockRepo)
		handler.ServeHTTP(rr, req)

		AssertStatusCode(t, rr, http.StatusBadRequest)

		expectedError := "json: cannot unmarshal"
		assert.Contains(t, rr.Body.String(), expectedError)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return status 500: failed to add record", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRecord := &models.EnergyRecord{
			Usage:  100,
			Device: "Laptop",
		}

		mockRepo.On("AddRecord", mockRecord).Return(fmt.Errorf("database_error"))

		body, _ := json.Marshal(mockRecord)
		req, rr := MakeRequest("POST", routeApi, body)

		handler := handlers.AddRecord(mockRepo)
		handler.ServeHTTP(rr, req)

		AssertStatusCode(t, rr, http.StatusInternalServerError)

		expectedError := "database_error\n"
		assert.Equal(t, expectedError, rr.Body.String())

		mockRepo.AssertExpectations(t)
	})
}

func TestGetRecords(t *testing.T) {
	const datePattern = "2006-01-02"
	const routeApi = "/api/records"

	t.Run("should get status 200: success get records", func(t *testing.T) {
		mockRepo := new(MockRepository)

		date1, _ := time.Parse(datePattern, "2023-12-31")
		date2, _ := time.Parse(datePattern, "2023-12-30")

		expectedRecords := []models.EnergyRecord{
			{ID: 1, Usage: 100, Device: "Air Conditioner", Date: date1},
			{ID: 2, Usage: 200, Device: "Refrigerator", Date: date2},
		}

		mockRepo.On("GetRecords").Return(expectedRecords, nil)

		handler := handlers.GetRecords(mockRepo)

		req, _ := http.NewRequest("GET", routeApi, nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var actualRecords []models.EnergyRecord
		err := json.NewDecoder(rr.Body).Decode(&actualRecords)
		assert.NoError(t, err)
		assert.Equal(t, expectedRecords, actualRecords)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return status 500: failed get records", func(t *testing.T) {
		mockRepo := new(MockRepository)

		mockRepo.On("GetRecords").Return([]models.EnergyRecord{}, fmt.Errorf("database_error"))

		handler := handlers.GetRecords(mockRepo)

		req, _ := http.NewRequest("GET", routeApi, nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		AssertStatusCode(t, rr, http.StatusInternalServerError)

		expectedError := "database_error\n"
		assert.Equal(t, expectedError, rr.Body.String())

		mockRepo.AssertExpectations(t)
	})
}

func TestGetRecordByID(t *testing.T) {
	const datePattern = "2006-01-02"

	t.Run("should get status 200: success get one record", func(t *testing.T) {
		mockRepo := new(MockRepository)

		date1, _ := time.Parse(datePattern, "2023-12-31")
		expectedRecords := &models.EnergyRecord{ID: 1, Usage: 100, Device: "Air Conditioner", Date: date1}

		mockRepo.On("GetByIdRecord", "1").Return(expectedRecords, nil)

		handler := handlers.GetByIdRecords(mockRepo)

		req, _ := http.NewRequest("GET", "/api/records/1", nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var actualRecords *models.EnergyRecord
		err := json.NewDecoder(rr.Body).Decode(&actualRecords)
		assert.NoError(t, err)
		assert.Equal(t, expectedRecords, actualRecords)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should get status 400: invalid param id", func(t *testing.T) {
		mockRepo := new(MockRepository)

		mockRepo.On("GetByIdRecord", "1").Return(&models.EnergyRecord{}, nil)

		handler := handlers.GetByIdRecords(mockRepo)

		req, _ := http.NewRequest("GET", "/api/records/salah", nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "salah"})
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedError := "invalid param id\n"
		assert.Equal(t, expectedError, rr.Body.String())
		mockRepo.AssertNotCalled(t, "GetByIdRecord")
	})

	t.Run("should return status 500: failed get one record", func(t *testing.T) {
		mockRepo := new(MockRepository)

		mockRepo.On("GetByIdRecord", "1").Return(&models.EnergyRecord{}, fmt.Errorf("database error"))

		handler := handlers.GetByIdRecords(mockRepo)

		req, _ := http.NewRequest("GET", "/api/records/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		AssertStatusCode(t, rr, http.StatusInternalServerError)

		expectedError := "database error\n"
		assert.Equal(t, expectedError, rr.Body.String())

		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteRecord(t *testing.T) {
	t.Run("should get status 204: success delete record", func(t *testing.T) {
		mockRepo := new(MockRepository)

		mockRepo.On("DeleteRecord", "1").Return(nil)

		handler := handlers.DeleteRecords(mockRepo)

		req, _ := http.NewRequest("DELETE", "/api/records/1", nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should get status 400: invalid param id", func(t *testing.T) {
		mockRepo := new(MockRepository)

		mockRepo.On("DeleteRecord", "1").Return(nil)

		handler := handlers.DeleteRecords(mockRepo)

		req, _ := http.NewRequest("GET", "/api/records/salah", nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "salah"})
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedError := "invalid param id\n"
		assert.Equal(t, expectedError, rr.Body.String())
		mockRepo.AssertNotCalled(t, "DeleteRecord")
	})

	t.Run("should return status 500: failed delete record", func(t *testing.T) {
		mockRepo := new(MockRepository)

		mockRepo.On("DeleteRecord", "1").Return(fmt.Errorf("database error"))

		handler := handlers.DeleteRecords(mockRepo)

		req, _ := http.NewRequest("GET", "/api/records/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		AssertStatusCode(t, rr, http.StatusInternalServerError)

		expectedError := "database error\n"
		assert.Equal(t, expectedError, rr.Body.String())

		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateRecord(t *testing.T) {
	t.Run("should get status 200: success update record", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRecord := &models.EnergyRecord{
			Usage:  100,
			Device: "Laptop",
		}

		mockRepo.On("UpdateRecord", mockRecord).Return(nil)

		handler := handlers.UpdateRecords(mockRepo)

		body, _ := json.Marshal(mockRecord)
		req, rr := MakeRequest("PUT", "/api/records/1", body)

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)


		var response models.EnergyRecord
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		assert.Equal(t, mockRecord, &response)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return status 400: invalid JSON body", func(t *testing.T) {
		mockRepo := new(MockRepository)
		invalidPayload := `{"Usage": "invalid_number", "Device": 123}`

		req, rr := MakeRequest("PUT", "/api/records/1", []byte(invalidPayload))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		handler := handlers.UpdateRecords(mockRepo)
		handler.ServeHTTP(rr, req)

		AssertStatusCode(t, rr, http.StatusBadRequest)

		expectedError := "json: cannot unmarshal"
		assert.Contains(t, rr.Body.String(), expectedError)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should get status 400: invalid request", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRecord := &models.EnergyRecord{}

		mockRepo.On("UpdateRecord", mockRecord).Return(nil)

		body, _ := json.Marshal(mockRecord)
		req, rr := MakeRequest("PUT", "/api/records/1", body)
		
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		
		handler := handlers.UpdateRecords(mockRepo)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedError := "usage is required and must be greater than 0\n"
		assert.Equal(t, expectedError, rr.Body.String())
		mockRepo.AssertNotCalled(t, "UpdateRecord")
	})

	t.Run("should get status 400: invalid param id", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRecord := &models.EnergyRecord{}

		mockRepo.On("UpdateRecord", mockRecord).Return(nil)

		body, _ := json.Marshal(mockRecord)
		req, rr := MakeRequest("PUT", "/api/records/salah", body)
		
		req = mux.SetURLVars(req, map[string]string{"id": "salah"})
		
		handler := handlers.UpdateRecords(mockRepo)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedError := "invalid param id\n"
		assert.Equal(t, expectedError, rr.Body.String())
		mockRepo.AssertNotCalled(t, "UpdateRecord")
	})

	t.Run("should return status 500: failed delete record", func(t *testing.T) {
		mockRepo := new(MockRepository)	
		mockRecord := &models.EnergyRecord{
			Usage:  100,
			Device: "Laptop",
		}
		
		mockRepo.On("UpdateRecord", mockRecord).Return(fmt.Errorf("database error"))

		body, _ := json.Marshal(mockRecord)
		req, rr := MakeRequest("PUT", "/api/records/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		
		handler := handlers.UpdateRecords(mockRepo)
		handler.ServeHTTP(rr, req)

		AssertStatusCode(t, rr, http.StatusInternalServerError)

		expectedError := "database error\n"
		assert.Equal(t, expectedError, rr.Body.String())

		mockRepo.AssertExpectations(t)
	})
}
