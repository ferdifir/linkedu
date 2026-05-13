package schedule

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

// Schedule handlers
func (h *Handler) GetSchedules(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	schedules, err := h.service.GetSchedules(*tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data jadwal"))
	}

	return c.JSON(shared.SuccessResponse(schedules, ""))
}

func (h *Handler) GetScheduleByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	schedule, err := h.service.GetScheduleByID(uint(id), *tenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Jadwal tidak ditemukan"))
	}

	return c.JSON(shared.SuccessResponse(schedule, ""))
}

type CreateScheduleInput struct {
	AcademicYearID uint   `json:"academic_year_id"`
	ClassroomID    uint   `json:"classroom_id"`
	SubjectID      uint   `json:"subject_id"`
	TeacherID      uint   `json:"teacher_id"`
	DayOfWeek      int    `json:"day_of_week"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
}

func (h *Handler) CreateSchedule(c *fiber.Ctx) error {
	var input CreateScheduleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	schedule, err := h.service.CreateSchedule(*tenantID, input.AcademicYearID, input.ClassroomID, input.SubjectID, input.TeacherID, input.DayOfWeek, input.StartTime, input.EndTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuat jadwal"))
	}

	return c.JSON(shared.SuccessResponse(schedule, "Jadwal berhasil dibuat"))
}

func (h *Handler) UpdateSchedule(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	var input CreateScheduleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	schedule, err := h.service.UpdateSchedule(uint(id), *tenantID, input.AcademicYearID, input.ClassroomID, input.SubjectID, input.TeacherID, input.DayOfWeek, input.StartTime, input.EndTime)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Jadwal tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal memperbarui jadwal"))
	}

	return c.JSON(shared.SuccessResponse(schedule, "Jadwal berhasil diperbarui"))
}

func (h *Handler) DeleteSchedule(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	if err := h.service.DeleteSchedule(uint(id), *tenantID); err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Jadwal tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menghapus jadwal"))
	}

	return c.JSON(shared.SuccessResponse(nil, "Jadwal berhasil dihapus"))
}

// Event handlers
func (h *Handler) GetEvents(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	events, err := h.service.GetEvents(*tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data event"))
	}

	return c.JSON(shared.SuccessResponse(events, ""))
}

func (h *Handler) GetEventByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	event, err := h.service.GetEventByID(uint(id), *tenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Event tidak ditemukan"))
	}

	return c.JSON(shared.SuccessResponse(event, ""))
}

type CreateEventInput struct {
	Name        string `json:"name"`
	EventType   string `json:"event_type"`
	Date        string `json:"date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description"`
}

func (h *Handler) CreateEvent(c *fiber.Ctx) error {
	var input CreateEventInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	event, err := h.service.CreateEvent(*tenantID, input.Name, input.EventType, input.Date, input.StartTime, input.EndTime, input.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuat event"))
	}

	return c.JSON(shared.SuccessResponse(event, "Event berhasil dibuat"))
}

func (h *Handler) UpdateEvent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	var input CreateEventInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	event, err := h.service.UpdateEvent(uint(id), *tenantID, input.Name, input.EventType, input.Date, input.StartTime, input.EndTime, input.Description)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Event tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal memperbarui event"))
	}

	return c.JSON(shared.SuccessResponse(event, "Event berhasil diperbarui"))
}

func (h *Handler) DeleteEvent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	if err := h.service.DeleteEvent(uint(id), *tenantID); err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Event tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menghapus event"))
	}

	return c.JSON(shared.SuccessResponse(nil, "Event berhasil dihapus"))
}
