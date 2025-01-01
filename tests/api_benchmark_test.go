package test

import (
	"daya-listrik-api/internal/handlers"
	"daya-listrik-api/internal/models"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func BenchmarkGetRecords(b *testing.B) {
	const datePattern = "2006-01-02"
	const routeApi = "/api/records"
	mockRepo := new(MockRepository)

	date1, _ := time.Parse(datePattern, "2023-12-31")
	date2, _ := time.Parse(datePattern, "2023-12-30")

	expectedRecords := []models.EnergyRecord{
		{ID: 1, Usage: 100, Device: "Air Conditioner", Date: date1},
		{ID: 2, Usage: 200, Device: "Refrigerator", Date: date2},
	}

	mockRepo.On("GetRecords").Return(expectedRecords, nil)

	handler := handlers.GetRecords(mockRepo)

	b.ResetTimer() // Mereset timer untuk memastikan hanya bagian pengujian yang dihitung
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", routeApi, nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			b.Errorf("Expected status 200, got %d", rr.Code)
		}

		var actualRecords []models.EnergyRecord
		err := json.NewDecoder(rr.Body).Decode(&actualRecords)
		if err != nil {
			b.Errorf("Error decoding response body: %v", err)
		}
		if len(actualRecords) != len(expectedRecords) {
			b.Errorf("Expected %d records, got %d", len(expectedRecords), len(actualRecords))
		}
	}

	mockRepo.AssertExpectations(b)
}

func BenchmarkAddRecord(b *testing.B) {
	const routeApi = "/api/records/add"
	mockRepo := new(MockRepository)
	mockRecord := &models.EnergyRecord{
		Usage:  100,
		Device: "Laptop",
	}

	mockRepo.On("AddRecord", mockRecord).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body, _ := json.Marshal(mockRecord)
		req, rr := MakeRequest("POST", routeApi, body)

		handler := handlers.AddRecord(mockRepo)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			b.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
		}

		var response models.EnergyRecord
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			b.Errorf("Error decoding response: %v", err)
		}

		if response != *mockRecord {
			b.Errorf("Expected record %+v, got %+v", mockRecord, response)
		}
	}

	mockRepo.AssertExpectations(b)
}

func BenchmarkGetByIdRecord(b *testing.B) {
	const datePattern = "2006-01-02"
	mockRepo := new(MockRepository)

	date1, _ := time.Parse(datePattern, "2023-12-31")
	expectedRecords := &models.EnergyRecord{ID: 1, Usage: 100, Device: "Air Conditioner", Date: date1}

	mockRepo.On("GetByIdRecord", "1").Return(expectedRecords, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler := handlers.GetByIdRecords(mockRepo)

		req, _ := http.NewRequest("GET", "/api/records/1", nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		b.StartTimer()

		handler.ServeHTTP(rr, req)

		b.StopTimer()

		var actualRecords *models.EnergyRecord
		err := json.NewDecoder(rr.Body).Decode(&actualRecords)
		if err != nil {
			b.Fatalf("Failed to decode response: %v", err)
		}

		if !assert.Equal(b, expectedRecords, actualRecords) {
			b.Fatalf("Expected %v but got %v", expectedRecords, actualRecords)
		}
	}

	mockRepo.AssertExpectations(b)
}

func BenchmarkDeleteRecord(b *testing.B) {
	mockRepo := new(MockRepository)

	mockRepo.On("DeleteRecord", "1").Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler := handlers.DeleteRecords(mockRepo)

		req, _ := http.NewRequest("DELETE", "/api/records/1", nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		b.StartTimer()

		handler.ServeHTTP(rr, req)

		b.StopTimer()

		if rr.Code != http.StatusNoContent {
			b.Fatalf("Expected status %d but got %d", http.StatusNoContent, rr.Code)
		}
	}

	mockRepo.AssertExpectations(b)
}

func BenchmarkUpdateRecord(b *testing.B) {
	mockRepo := new(MockRepository)
	mockRecord := &models.EnergyRecord{
		Usage:  100,
		Device: "Laptop",
	}
	mockRepo.On("UpdateRecord", mockRecord).Return(nil)

	// Start benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler := handlers.UpdateRecords(mockRepo)

		body, err := json.Marshal(mockRecord)
		if err != nil {
			b.Fatalf("Error marshalling mock record: %v", err)
		}

		req, rr := MakeRequest("PUT", "/api/records/1", body)

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		b.StartTimer()

		handler.ServeHTTP(rr, req)

		b.StopTimer()

		if rr.Code != http.StatusOK {
			b.Fatalf("Expected status %d but got %d", http.StatusOK, rr.Code)
		}

		var response models.EnergyRecord
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			b.Fatalf("Error decoding response: %v", err)
		}

		if !assert.Equal(b, mockRecord, &response) {
			b.Fatalf("Expected record: %v but got: %v", mockRecord, &response)
		}
	}

	mockRepo.AssertExpectations(b)
}
