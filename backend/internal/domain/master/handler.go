package master

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

// Classroom handlers
func (h *Handler) GetClassrooms(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	classrooms, err := h.service.GetClassrooms(*tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data kelas"))
	}

	return c.JSON(shared.SuccessResponse(classrooms, ""))
}

func (h *Handler) GetClassroomByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	classroom, err := h.service.GetClassroomByID(uint(id), *tenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Kelas tidak ditemukan"))
	}

	return c.JSON(shared.SuccessResponse(classroom, ""))
}

type CreateClassroomInput struct {
	Name       string `json:"name"`
	GradeLevel string `json:"grade_level"`
}

func (h *Handler) CreateClassroom(c *fiber.Ctx) error {
	var input CreateClassroomInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	classroom, err := h.service.CreateClassroom(*tenantID, input.Name, input.GradeLevel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuat kelas"))
	}

	return c.JSON(shared.SuccessResponse(classroom, "Kelas berhasil dibuat"))
}

func (h *Handler) UpdateClassroom(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	var input CreateClassroomInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	classroom, err := h.service.UpdateClassroom(uint(id), *tenantID, input.Name, input.GradeLevel)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Kelas tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal memperbarui kelas"))
	}

	return c.JSON(shared.SuccessResponse(classroom, "Kelas berhasil diperbarui"))
}

func (h *Handler) DeleteClassroom(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	if err := h.service.DeleteClassroom(uint(id), *tenantID); err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Kelas tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menghapus kelas"))
	}

	return c.JSON(shared.SuccessResponse(nil, "Kelas berhasil dihapus"))
}

// Subject handlers
func (h *Handler) GetSubjects(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	subjects, err := h.service.GetSubjects(*tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data mata pelajaran"))
	}

	return c.JSON(shared.SuccessResponse(subjects, ""))
}

type CreateSubjectInput struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (h *Handler) CreateSubject(c *fiber.Ctx) error {
	var input CreateSubjectInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	subject, err := h.service.CreateSubject(*tenantID, input.Name, input.Code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuat mata pelajaran"))
	}

	return c.JSON(shared.SuccessResponse(subject, "Mata pelajaran berhasil dibuat"))
}

func (h *Handler) UpdateSubject(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	var input CreateSubjectInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	subject, err := h.service.UpdateSubject(uint(id), *tenantID, input.Name, input.Code)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Mata pelajaran tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal memperbarui mata pelajaran"))
	}

	return c.JSON(shared.SuccessResponse(subject, "Mata pelajaran berhasil diperbarui"))
}

func (h *Handler) DeleteSubject(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	if err := h.service.DeleteSubject(uint(id), *tenantID); err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Mata pelajaran tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menghapus mata pelajaran"))
	}

	return c.JSON(shared.SuccessResponse(nil, "Mata pelajaran berhasil dihapus"))
}

// Teacher handlers
func (h *Handler) GetTeachers(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	teachers, err := h.service.GetTeachers(*tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data guru"))
	}

	return c.JSON(shared.SuccessResponse(teachers, ""))
}

type CreateTeacherInput struct {
	UserID uint   `json:"user_id"`
	NIP    string `json:"nip"`
	Phone  string `json:"phone"`
}

func (h *Handler) CreateTeacher(c *fiber.Ctx) error {
	var input CreateTeacherInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	teacher, err := h.service.CreateTeacher(*tenantID, input.UserID, input.NIP, input.Phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuat guru"))
	}

	return c.JSON(shared.SuccessResponse(teacher, "Guru berhasil dibuat"))
}

func (h *Handler) UpdateTeacher(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	var input CreateTeacherInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	teacher, err := h.service.UpdateTeacher(uint(id), *tenantID, input.NIP, input.Phone)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Guru tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal memperbarui guru"))
	}

	return c.JSON(shared.SuccessResponse(teacher, "Guru berhasil diperbarui"))
}

func (h *Handler) DeleteTeacher(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	if err := h.service.DeleteTeacher(uint(id), *tenantID); err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Guru tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menghapus guru"))
	}

	return c.JSON(shared.SuccessResponse(nil, "Guru berhasil dihapus"))
}

// Student handlers
func (h *Handler) GetStudents(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	students, err := h.service.GetStudents(*tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data siswa"))
	}

	return c.JSON(shared.SuccessResponse(students, ""))
}

type CreateStudentInput struct {
	UserID      uint   `json:"user_id"`
	ClassroomID uint   `json:"classroom_id"`
	NIS         string `json:"nis"`
	RFIDUID     string `json:"rfid_uid"`
}

func (h *Handler) CreateStudent(c *fiber.Ctx) error {
	var input CreateStudentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	student, err := h.service.CreateStudent(*tenantID, input.UserID, input.ClassroomID, input.NIS, input.RFIDUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuat siswa"))
	}

	return c.JSON(shared.SuccessResponse(student, "Siswa berhasil dibuat"))
}

func (h *Handler) UpdateStudent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	var input CreateStudentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	student, err := h.service.UpdateStudent(uint(id), *tenantID, input.ClassroomID, input.NIS, input.RFIDUID)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Siswa tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal memperbarui siswa"))
	}

	return c.JSON(shared.SuccessResponse(student, "Siswa berhasil diperbarui"))
}

func (h *Handler) DeleteStudent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	if err := h.service.DeleteStudent(uint(id), *tenantID); err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Siswa tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menghapus siswa"))
	}

	return c.JSON(shared.SuccessResponse(nil, "Siswa berhasil dihapus"))
}

// Parent handlers
func (h *Handler) GetParents(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	parents, err := h.service.GetParents(*tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data orang tua"))
	}

	return c.JSON(shared.SuccessResponse(parents, ""))
}

type CreateParentInput struct {
	UserID    uint `json:"user_id"`
	StudentID uint `json:"student_id"`
}

func (h *Handler) CreateParent(c *fiber.Ctx) error {
	var input CreateParentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	parent, err := h.service.CreateParent(*tenantID, input.UserID, input.StudentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuat orang tua"))
	}

	return c.JSON(shared.SuccessResponse(parent, "Orang tua berhasil dibuat"))
}

func (h *Handler) UpdateParent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	var input CreateParentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	parent, err := h.service.UpdateParent(uint(id), *tenantID, input.StudentID)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Orang tua tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal memperbarui orang tua"))
	}

	return c.JSON(shared.SuccessResponse(parent, "Orang tua berhasil diperbarui"))
}

func (h *Handler) DeleteParent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	if err := h.service.DeleteParent(uint(id), *tenantID); err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Orang tua tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menghapus orang tua"))
	}

	return c.JSON(shared.SuccessResponse(nil, "Orang tua berhasil dihapus"))
}

// Academic Year handlers
func (h *Handler) GetAcademicYears(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	academicYears, err := h.service.GetAcademicYears(*tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mengambil data tahun ajaran"))
	}

	return c.JSON(shared.SuccessResponse(academicYears, ""))
}

type CreateAcademicYearInput struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	IsActive  bool   `json:"is_active"`
}

func (h *Handler) CreateAcademicYear(c *fiber.Ctx) error {
	var input CreateAcademicYearInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	academicYear, err := h.service.CreateAcademicYear(*tenantID, input.Name, input.StartDate, input.EndDate, input.IsActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuat tahun ajaran"))
	}

	return c.JSON(shared.SuccessResponse(academicYear, "Tahun ajaran berhasil dibuat"))
}

func (h *Handler) UpdateAcademicYear(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	var input CreateAcademicYearInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	academicYear, err := h.service.UpdateAcademicYear(uint(id), *tenantID, input.Name, input.StartDate, input.EndDate, input.IsActive)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Tahun ajaran tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal memperbarui tahun ajaran"))
	}

	return c.JSON(shared.SuccessResponse(academicYear, "Tahun ajaran berhasil diperbarui"))
}

func (h *Handler) DeleteAcademicYear(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("ID tidak valid"))
	}

	tenantID := c.Locals("tenant_id").(*uint)
	if tenantID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Tenant ID tidak ditemukan"))
	}

	if err := h.service.DeleteAcademicYear(uint(id), *tenantID); err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Tahun ajaran tidak ditemukan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal menghapus tahun ajaran"))
	}

	return c.JSON(shared.SuccessResponse(nil, "Tahun ajaran berhasil dihapus"))
}
