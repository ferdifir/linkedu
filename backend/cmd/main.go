package main

import (
	"log"

	"linkedu/internal/config"
	"linkedu/internal/database"
	"linkedu/internal/domain"
	"linkedu/internal/domain/attendance"
	"linkedu/internal/domain/auth"
	"linkedu/internal/domain/tenant"
	"linkedu/internal/domain/user"
	"linkedu/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := database.Init(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto migrate models
	if err := database.GetDB().AutoMigrate(
		&domain.Tenant{},
		&domain.User{},
		&domain.AcademicYear{},
		&domain.Classroom{},
		&domain.Subject{},
		&domain.Teacher{},
		&domain.Student{},
		&domain.Parent{},
		&domain.Schedule{},
		&domain.Event{},
		&domain.AttendanceSession{},
		&domain.AttendanceRecord{},
		&domain.Permit{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize JWT
	middleware.InitJWT(cfg.JWTSecret)

	// Initialize repositories
	tenantRepo := tenant.NewRepository()
	userRepo := user.NewRepository()
	attendanceRepo := attendance.NewRepository()

	// Initialize services
	tenantService := tenant.NewService(tenantRepo)
	userService := user.NewService(userRepo)
	attendanceService := attendance.NewService(attendanceRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(tenantService, userService)
	attendanceHandler := attendance.NewHandler(attendanceService)

	// Create Fiber app
	app := fiber.New()

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	// Public routes
	api := app.Group("/api/v1")
	authGroup := api.Group("/auth")
	authGroup.Post("/register-tenant", authHandler.RegisterTenant)
	authGroup.Post("/login", authHandler.Login)

	// Protected routes
	protected := api.Group("", middleware.AuthMiddleware())

	// Attendance routes (teacher)
	attendanceGroup := protected.Group("/sessions")
	attendanceGroup.Post("", attendanceHandler.OpenSession)
	attendanceGroup.Patch("/:id/close", attendanceHandler.CloseSession)
	attendanceGroup.Post("/:id/tap", attendanceHandler.TapStudent)
	attendanceGroup.Get("/:id/records", attendanceHandler.GetSessionRecords)

	// Start server
	log.Printf("Server starting on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
