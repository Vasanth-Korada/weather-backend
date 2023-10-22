package repository

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	if app == nil {
		log.Fatal("Application instance cannot be nil")
	}

	if r.DB == nil {
		log.Fatal("Database instance cannot be nil")
	}

	api := app.Group("/api")

	api.Use(r.loggingMiddleware)

	api.Post("/create_weather", r.ensureDBConnection(r.CreateWeather))
	api.Delete("/delete_weather/:id", r.ensureDBConnection(r.DeleteWeather))
	api.Put("/update_weather/:id", r.ensureDBConnection(r.UpdateWeather))
	api.Get("/history", r.ensureDBConnection(r.GetHistory))
	api.Post("/login", r.ensureDBConnection(r.handleLogin))
	api.Post("/register", r.ensureDBConnection(r.handleRegister))
	api.Post("/logout", r.ensureDBConnection(r.handleLogout))
}

func (r *Repository) loggingMiddleware(c *fiber.Ctx) error {
	log.Printf("API Request made to %s", c.OriginalURL())
	return c.Next()
}

func (r *Repository) ensureDBConnection(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if r.DB == nil {
			log.Println("Database connection is not initialized")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}
		return next(c)
	}
}
