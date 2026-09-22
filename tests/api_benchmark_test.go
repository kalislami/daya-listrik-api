package test

import (
	"bytes"
	"daya-listrik-api/internal/handlers"
	"daya-listrik-api/internal/models"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	HeaderContentType     = "Content-Type"
	HeaderApplicationJSON = "application/json"
	pathRecords           = "/api/records"
	pathRecordsId         = "/api/records/:id"
	pathRecordsIdVal      = "/api/records/1"
)

func BenchmarkGetRecords(b *testing.B) {
	const datePattern = "2006-01-02"
	mockRepo := new(MockRepository)

	date1, _ := time.Parse(datePattern, "2023-12-31")
	date2, _ := time.Parse(datePattern, "2023-12-30")

	expectedRecords := []models.EnergyRecord{
		{ID: 1, Usage: 100, Device: "Air Conditioner", Date: date1},
		{ID: 2, Usage: 200, Device: "Refrigerator", Date: date2},
	}

	mockRepo.On("GetRecords", mock.Anything).Return(expectedRecords, nil)

	handler := &handlers.EnergyRecordHandler{Repo: mockRepo}
	app := fiber.New()
	app.Get(pathRecords, handler.GetRecords)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", pathRecords, nil)
		resp, _ := app.Test(req)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	}
	mockRepo.AssertExpectations(b)
}

func BenchmarkAddRecord(b *testing.B) {
	mockRepo := new(MockRepository)
	mockRecord := &models.EnergyRecord{
		Usage:    100,
		Device:   "Laptop",
		Duration: 1,
	}

	mockRepo.On("AddRecord", mock.Anything, mockRecord).Return(nil)

	handler := &handlers.EnergyRecordHandler{Repo: mockRepo}
	app := fiber.New()
	app.Post(pathRecords, handler.AddRecord)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body := []byte(`{"usage":100,"device":"Laptop","duration":1}`)
		req := httptest.NewRequest(http.MethodPost, pathRecords, bytes.NewReader(body))
		req.Header.Set(HeaderContentType, HeaderApplicationJSON)

		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusCreated {
			b.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		var response models.EnergyRecord
		err := json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			b.Errorf("Error decoding response: %v", err)
		}

		if !assert.Equal(b, mockRecord.Usage, response.Usage) ||
			!assert.Equal(b, mockRecord.Device, response.Device) {
			b.Errorf("Expected record %+v, got %+v", mockRecord, response)
		}
		resp.Body.Close()
	}

	mockRepo.AssertExpectations(b)
}

func BenchmarkGetByIdRecord(b *testing.B) {
	mockRepo := new(MockRepository)

	expectedRecord := &models.EnergyRecord{ID: 1, Usage: 100, Device: "Air Conditioner", Date: time.Now()}
	mockRepo.On("GetByIdRecord", mock.Anything, 1).Return(expectedRecord, nil)

	handler := &handlers.EnergyRecordHandler{Repo: mockRepo}
	app := fiber.New()
	app.Get(pathRecordsId, handler.GetByIdRecord)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, pathRecordsIdVal, nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var actualRecord models.EnergyRecord
		err := json.NewDecoder(resp.Body).Decode(&actualRecord)
		if err != nil {
			b.Fatalf("Failed to decode response: %v", err)
		}

		if !assert.Equal(b, expectedRecord.ID, actualRecord.ID) {
			b.Fatalf("Expected ID %d but got %d", expectedRecord.ID, actualRecord.ID)
		}
		resp.Body.Close()
	}

	mockRepo.AssertExpectations(b)
}

func BenchmarkDeleteRecord(b *testing.B) {
	mockRepo := new(MockRepository)
	mockRepo.On("DeleteRecord", mock.Anything, 1).Return(nil)

	handler := &handlers.EnergyRecordHandler{Repo: mockRepo}
	app := fiber.New()
	app.Delete(pathRecordsId, handler.DeleteRecord)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodDelete, pathRecordsIdVal, nil)
		resp, _ := app.Test(req)
		resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			b.Fatalf("Expected status %d but got %d", http.StatusNoContent, resp.StatusCode)
		}
	}

	mockRepo.AssertExpectations(b)
}

func BenchmarkUpdateRecord(b *testing.B) {
	mockRepo := new(MockRepository)
	mockRecord := &models.EnergyRecord{ID: 1, Usage: 100, Device: "Laptop", Duration: 1}
	mockRepo.On("UpdateRecord", mock.Anything, mockRecord).Return(nil)

	handler := &handlers.EnergyRecordHandler{Repo: mockRepo}
	app := fiber.New()
	app.Put(pathRecordsId, handler.UpdateRecord)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body := []byte(`{"usage":100,"device":"Laptop","duration":1}`)

		req := httptest.NewRequest(http.MethodPut, pathRecordsIdVal, bytes.NewReader(body))
		req.Header.Set(HeaderContentType, HeaderApplicationJSON)

		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status %d but got %d", http.StatusOK, resp.StatusCode)
		}

		var response models.EnergyRecord
		err := json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			b.Fatalf("Error decoding response: %v", err)
		}

		if !assert.Equal(b, mockRecord.ID, response.ID) {
			b.Fatalf("Expected record ID %v but got %v", mockRecord.ID, response.ID)
		}
		resp.Body.Close()
	}

	mockRepo.AssertExpectations(b)
}
