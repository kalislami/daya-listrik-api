package handlers

import (
	"daya-listrik-api/internal/models"
	"daya-listrik-api/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func InitializeRoutes(r *mux.Router, repo repository.EnergyRecordRepositoryInterface) {
	const routeApiRecordsAdd = "/api/records/add"
	const routeApiRecord = "/api/records"
	const routeApiRecordsId = "/api/records/{id}"

	r.HandleFunc(routeApiRecord, GetRecords(repo)).Methods("GET")
	r.HandleFunc(routeApiRecordsAdd, AddRecord(repo)).Methods("POST")
	r.HandleFunc(routeApiRecordsId, DeleteRecords(repo)).Methods("DELETE")
	r.HandleFunc(routeApiRecordsId, UpdateRecords(repo)).Methods("PUT")
	r.HandleFunc(routeApiRecordsId, GetByIdRecords(repo)).Methods("GET")
}

func validateEnergyRecord(record *models.EnergyRecord) error {
	if record.Usage <= 0 {
		return fmt.Errorf("usage is required and must be greater than 0")
	}
	if record.Device == "" {
		return fmt.Errorf("device is required")
	}
	return nil
}

func validateParamId(r *http.Request) (string, error) {
	id := strings.TrimSpace(mux.Vars(r)["id"])

	if _, err := strconv.Atoi(id); err != nil {
		return "", fmt.Errorf("invalid param id")
	}

	return id, nil
}

func AddRecord(repo repository.EnergyRecordRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var record models.EnergyRecord
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			log.Printf("Invalid JSON: %v", err)
			http.Error(w, "Input tidak valid. Pastikan semua nilai benar.", http.StatusBadRequest)
			return
		}

		if err := validateEnergyRecord(&record); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repo.AddRecord(&record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(record)
	}
}

func GetRecords(repo repository.EnergyRecordRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records, err := repo.GetRecords()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(records)
	}
}

func DeleteRecords(repo repository.EnergyRecordRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := validateParamId(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repo.DeleteRecord(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateRecords(repo repository.EnergyRecordRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := validateParamId(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idInt, _ := strconv.Atoi(id)

		var record models.EnergyRecord
		record.ID = idInt
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateEnergyRecord(&record); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repo.UpdateRecord(&record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(record)
	}
}

func GetByIdRecords(repo repository.EnergyRecordRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := validateParamId(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		record, err := repo.GetByIdRecord(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(record)
	}
}
