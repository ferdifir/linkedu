package auth

import (
	"linkedu/internal/domain/tenant"
	"linkedu/internal/domain/user"
	"linkedu/internal/middleware"
	"linkedu/internal/shared"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	tenantService *tenant.Service
	userService   *user.Service
}

func NewHandler(tenantSvc *tenant.Service, userSvc *user.Service) *Handler {
	return &Handler{
		tenantService: tenantSvc,
		userService:   userSvc,
	}
}

type RegisterTenantInput struct {
	Name       string `json:"name"`
	SchoolName string `json:"school_name"`
	Email      string `json:"email"`
	Subdomain  string `json:"subdomain"`
	Password   string `json:"password"`
}

func (h *Handler) RegisterTenant(c *fiber.Ctx) error {
	var input RegisterTenantInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	tenant, user, err := h.tenantService.RegisterTenant(tenant.RegisterTenantInput(input))
	if err != nil {
		if err == tenant.ErrEmailExists {
			return c.Status(fiber.StatusConflict).JSON(shared.ErrorResponse("Email sudah terdaftar"))
		}
		if err == tenant.ErrSubdomainExists {
			return c.Status(fiber.StatusConflict).JSON(shared.ErrorResponse("Subdomain sudah digunakan"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal mendaftar"))
	}

	token, err := middleware.GenerateToken(user.ID, user.TenantID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuat token"))
	}

	return c.Status(fiber.StatusCreated).JSON(shared.SuccessResponse(fiber.Map{
		"tenant":       tenant,
		"user":         user,
		"access_token": token,
	}, "Pendaftaran berhasil"))
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse("Request tidak valid"))
	}

	output, err := h.userService.Login(user.LoginInput(input))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse(err.Error()))
	}

	token, err := middleware.GenerateToken(output.User.ID, output.User.TenantID, output.User.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse("Gagal membuat token"))
	}

	return c.JSON(shared.SuccessResponse(fiber.Map{
		"user":         output.User,
		"access_token": token,
	}, "Login berhasil"))
}
