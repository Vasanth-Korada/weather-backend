package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Vasanth-Korada/weather-backend/storage"

	"github.com/Vasanth-Korada/weather-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Weather struct {
	Temperature string `json:"author"`
	City        string `json:"city"`
	Lat         string `json:"lat"`
	Long        string `json:"lng"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateWeather(context *fiber.Ctx) error {
	weather := Weather{}

	err := context.BodyParser(&weather)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&weather).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create a weather record"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "weather record has been added"})
	return nil
}

func (r *Repository) GetHistory(context *fiber.Ctx) error {
	weatherModels := &[]models.Weather{}
	err := r.DB.Find(weatherModels).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get search history"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "history fetched successfully", "data": weatherModels})
	return nil

}

func (r *Repository) DeleteWeather(context *fiber.Ctx) error {
	weatherModel := models.Weather{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	err := r.DB.Delete(weatherModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete weather",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "weather deleted successfully",
	})
	return nil
}

func (r *Repository) UpdateWeather(context *fiber.Ctx) error {
	id := context.Params("id")
	lat := context.Params("lat")
	lng := context.Params("lat")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	err := r.DB.Where("id = ?", id).Update("lat", lat).Update("lng", lng)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not update weather",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "weather updated successfully",
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_weather", r.CreateWeather)
	api.Delete("/delete_weather/:id", r.DeleteWeather)
	app.Put("/update_weather/:id", r.UpdateWeather)
	app.Get("/history", r.GetHistory)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("Could not load the database")
	}

	err = models.MigrateWeather(db)

	if err != nil {
		log.Fatal("Could not migrate database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
