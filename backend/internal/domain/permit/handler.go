package permit

import (
	"linkedu/internal/shared"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Get all permits (for admin)
func (h *Handler) GetPermits(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	permits, err := h.service.GetPermits(*tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data izin"))
	}

	return c.JSON(shared.SuccessResponse(permits, ""))
}

// Get permit by ID
func (h *Handler) GetPermitByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	permit, err := h.service.GetPermitByID(uint(id), *tenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Izin tidak ditemukan"))
	}

	return c.JSON(shared.SuccessResponse(permit, ""))
}

// Get permits for student
func (h *Handler) GetPermitsByStudent(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("student_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID siswa tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	permits, err := h.service.GetPermitsByStudent(uint(studentID), *tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data izin"))
	}

	return c.JSON(shared.SuccessResponse(permits, ""))
}

// Get permits for parent
func (h *Handler) GetPermitsByParent(c *fiber.Ctx) error {
	parentID, err := strconv.ParseUint(c.Params("parent_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID orang tua tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	permits, err := h.service.GetPermitsByParent(uint(parentID), *tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data izin"))
	}

	return c.JSON(shared.SuccessResponse(permits, ""))
}

// Get pending permits (for admin)
func (h *Handler) GetPendingPermits(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	permits, err := h.service.GetPermitsByStatus("pending", *tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data izin"))
	}

	return c.JSON(shared.SuccessResponse(permits, ""))
}

type CreatePermitInput struct {
	StudentID uint   `json:"student_id"`
	SessionID *uint  `json:"session_id"`
	Date      string `json:"date"`
	Reason    string `json:"reason"`
}

func (h *Handler) CreatePermit(c *fiber.Ctx) error {
	var input CreatePermitInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	parentID := c.Locals("user_id").(uint)

	// Handle file upload
	var attachmentPath string
	file, err := c.FormFile("attachment")
	if err == nil {
		// Create storage directory if not exists
		storageDir := "./storage/permits"
		if err := os.MkdirAll(storageDir, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menyimpan file"))
		}

		// Generate unique filename
		filename := time.Now().Format("20060102150405") + "_" + file.Filename
		filePath := filepath.Join(storageDir, filename)

		// Save file
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menyimpan file"))
		}

		attachmentPath = filePath
	}

	permit, err := h.service.CreatePermit(*tenantID, input.StudentID, parentID, input.SessionID, input.Date, input.Reason, attachmentPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengajukan izin"))
	}

	return c.JSON(shared.SuccessResponse(permit, "Izin berhasil diajukan"))
}

type ReviewPermitInput struct {
	Status string `json:"status"` // approved or rejected
	Notes  string `json:"notes"`
}

func (h *Handler) ReviewPermit(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	var input ReviewPermitInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	reviewedBy := c.Locals("user_id").(uint)

	permit, err := h.service.ReviewPermit(uint(id), *tenantID, reviewedBy, input.Status, input.Notes)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Izin tidak ditemukan"))
		}
		if err == ErrInvalidStatus {
			return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Status tidak valid"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal memproses izin"))
	}

	return c.JSON(shared.SuccessResponse(permit, "Izin berhasil diproses"))
}

// Download attachment
func (h *Handler) DownloadAttachment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	permit, err := h.service.GetPermitByID(uint(id), *tenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Izin tidak ditemukan"))
	}

	if permit.AttachmentPath == "" {
		return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("File tidak ditemukan"))
	}

	return c.Download(permit.AttachmentPath)
}
