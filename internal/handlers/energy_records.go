package handlers

import (
	"bytes"
	"daya-listrik-api/internal/models"
	"daya-listrik-api/internal/repository"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
)

type EnergyRecordHandler struct {
	Repo repository.EnergyRecordStore
}

type recordRequest struct {
	Device   string  `json:"device"`
	Usage    float64 `json:"usage"`
	Duration float64 `json:"duration"`
}

type apiError struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func InitializeRoutes(app *fiber.App, handler *EnergyRecordHandler) {
	app.Get("/api/records", handler.GetRecords)
	app.Post("/api/records", handler.AddRecord)
	app.Get("/api/records/:id", handler.GetByIdRecord)
	app.Put("/api/records/:id", handler.UpdateRecord)
	app.Delete("/api/records/:id", handler.DeleteRecord)
}

func respondError(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(apiError{Error: errorDetail{Code: code, Message: message}})
}

func respondRepositoryError(c *fiber.Ctx, err error) error {
	if errors.Is(err, repository.ErrRecordNotFound) {
		return respondError(c, fiber.StatusNotFound, "RECORD_NOT_FOUND", "Energy record not found")
	}
	log.Printf("repository error: %v", err)
	return respondError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error")
}

func parseID(c *fiber.Ctx) (int, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid ID")
	}
	return id, nil
}

func parseRecord(c *fiber.Ctx) (recordRequest, error) {
	var request recordRequest
	contentType, _, err := mime.ParseMediaType(c.Get(fiber.HeaderContentType))
	if err != nil || contentType != fiber.MIMEApplicationJSON {
		return recordRequest{}, errors.New("expected application/json")
	}
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return recordRequest{}, err
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return recordRequest{}, errors.New("multiple JSON values")
	}
	request.Device = strings.TrimSpace(request.Device)
	if request.Device == "" || utf8.RuneCountInString(request.Device) > 100 || request.Usage <= 0 || request.Duration <= 0 {
		return recordRequest{}, errors.New("invalid record fields")
	}
	return request, nil
}

func (h *EnergyRecordHandler) AddRecord(c *fiber.Ctx) error {
	request, err := parseRecord(c)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Invalid request")
	}
	record := models.EnergyRecord{Device: request.Device, Usage: request.Usage, Duration: request.Duration}
	if err := h.Repo.AddRecord(c.Context(), &record); err != nil {
		return respondRepositoryError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(record)
}

func (h *EnergyRecordHandler) GetRecords(c *fiber.Ctx) error {
	records, err := h.Repo.GetRecords(c.Context())
	if err != nil {
		return respondRepositoryError(c, err)
	}
	return c.JSON(records)
}

func (h *EnergyRecordHandler) GetByIdRecord(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "INVALID_ID", "ID must be a positive integer")
	}
	record, err := h.Repo.GetByIdRecord(c.Context(), id)
	if err != nil {
		return respondRepositoryError(c, err)
	}
	return c.JSON(record)
}

func (h *EnergyRecordHandler) UpdateRecord(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "INVALID_ID", "ID must be a positive integer")
	}
	request, err := parseRecord(c)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Invalid request")
	}
	record := models.EnergyRecord{ID: id, Device: request.Device, Usage: request.Usage, Duration: request.Duration}
	if err := h.Repo.UpdateRecord(c.Context(), &record); err != nil {
		return respondRepositoryError(c, err)
	}
	return c.JSON(record)
}

func (h *EnergyRecordHandler) DeleteRecord(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "INVALID_ID", "ID must be a positive integer")
	}
	if err := h.Repo.DeleteRecord(c.Context(), id); err != nil {
		return respondRepositoryError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
