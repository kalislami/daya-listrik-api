package handlers

import (
	"daya-listrik-api/internal/models"
	"daya-listrik-api/internal/repository"
	"daya-listrik-api/internal/utils"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type EnergyRecordHandler struct {
	Repo repository.EnergyRecordRepositoryInterface
}

func InitializeRoutes(app *fiber.App, handler *EnergyRecordHandler) {
	const routeApiRecordsAdd = "/api/records/add"
	const routeApiRecord = "/api/records"
	const routeApiRecordId = "/api/records/:id"

	app.Get(routeApiRecord, handler.GetRecords)
	app.Post(routeApiRecordsAdd, handler.AddRecord)
	app.Delete(routeApiRecordId, handler.DeleteRecord)
	app.Put(routeApiRecordId, handler.UpdateRecord)
	app.Get(routeApiRecordId, handler.GetByIdRecord)
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

func validateParamId(c *fiber.Ctx) (string, error) {
	id := strings.TrimSpace(c.Params("id"))

	if _, err := strconv.Atoi(id); err != nil {
		return "", fmt.Errorf("invalid param id")
	}
	return id, nil
}

func (h *EnergyRecordHandler) AddRecord(c *fiber.Ctx) error {
	var record models.EnergyRecord
	if err := c.BodyParser(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Input tidak valid. Pastikan semua nilai benar.",
		})
	}

	if err := validateEnergyRecord(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Repo.AddRecord(c.Context(), &record); err != nil {
		log.Printf("failed to insert record: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save record",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(record)
}

func (h *EnergyRecordHandler) GetRecords(c *fiber.Ctx) error {
	records, err := h.Repo.GetRecords(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

func (h *EnergyRecordHandler) DeleteRecord(c *fiber.Ctx) error {
	id, err := validateParamId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Repo.DeleteRecord(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *EnergyRecordHandler) UpdateRecord(c *fiber.Ctx) error {
	id, err := validateParamId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	idInt, _ := strconv.Atoi(id)

	var record models.EnergyRecord
	record.ID = idInt
	if err := c.BodyParser(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validateEnergyRecord(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Repo.UpdateRecord(c.Context(), &record); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(record)
}

func (h *EnergyRecordHandler) GetByIdRecord(c *fiber.Ctx) error {
	id, err := validateParamId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	record, err := h.Repo.GetByIdRecord(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(record)
}
