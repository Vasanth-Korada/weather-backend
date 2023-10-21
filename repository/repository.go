package repository

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_weather", r.CreateWeather)
	api.Delete("/delete_weather/:id", r.DeleteWeather)
	api.Put("/update_weather/:id", r.UpdateWeather)
	api.Get("/history", r.GetHistory)
	api.Post("/login", r.handleLogin)
	api.Post("/register", r.handleRegister)
	api.Post("/logout", r.handleLogout)
}
