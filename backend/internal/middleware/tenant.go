package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenant_id")
		if tenantID == nil {
			return c.Next()
		}
		c.Locals("tenant_id", tenantID)
		return c.Next()
	}
}
