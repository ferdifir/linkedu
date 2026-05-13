package attendance

import (
	"linkedu/internal/shared"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) OpenSession(c *fiber.Ctx) error {
	var input OpenSessionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	teacherID := c.Locals("user_id").(uint)
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	session, err := h.service.OpenSession(OpenSessionInput{
		ScheduleID: input.ScheduleID,
		EventID:    input.EventID,
		TeacherID:  teacherID,
		TenantID:   *tenantID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuka sesi presensi"))
	}

	return c.JSON(shared.SuccessResponse(session, "Sesi presensi dibuka"))
}

func (h *Handler) CloseSession(c *fiber.Ctx) error {
	sessionID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID sesi tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	if err := h.service.CloseSession(uint(sessionID), *tenantID); err != nil {
		if err == ErrSessionAlreadyClosed {
			return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse(err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menutup sesi presensi"))
	}

	return c.JSON(shared.SuccessResponse(nil, "Sesi presensi ditutup"))
}

func (h *Handler) TapStudent(c *fiber.Ctx) error {
	var input TapInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	sessionID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID sesi tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	record, student, err := h.service.TapStudent(uint(sessionID), input.RFIDUID, *tenantID)
	if err != nil {
		if err == ErrStudentNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Kartu tidak dikenal"))
		}
		if err == ErrSessionClosed {
			return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Sesi presensi sudah ditutup"))
		}
		if err == ErrSessionNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Sesi presensi tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mencatat presensi"))
	}

	return c.JSON(shared.SuccessResponse(fiber.Map{
		"record":  record,
		"student": student,
	}, "Presensi berhasil"))
}

func (h *Handler) GetSessionRecords(c *fiber.Ctx) error {
	sessionID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID sesi tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	records, err := h.service.repo.GetRecordsBySession(uint(sessionID), *tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data presensi"))
	}

	return c.JSON(shared.SuccessResponse(records, ""))
}
