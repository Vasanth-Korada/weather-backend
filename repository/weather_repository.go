package repository

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Vasanth-Korada/weather-backend/models"
	"github.com/gofiber/fiber/v2"
)

type Weather struct {
	Temperature string `json:"temp"`
	City        string `json:"city"`
	Lat         string `json:"lat"`
	Long        string `json:"lng"`
	UserID      int    `json:"userId"`
}

func (r *Repository) CreateWeather(context *fiber.Ctx) error {
	userID, ok := context.Locals("user").(int)
	if !ok {
		log.Println("Failed to extract userID from context")
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid user ID"})
	}

	var weather Weather
	if err := context.BodyParser(&weather); err != nil {
		log.Println("Error parsing request body:", err)
		return context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Request failed"})
	}

	weather.UserID = userID

	if err := r.DB.Create(&weather).Error; err != nil {
		log.Println("Error creating weather record:", err)
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not create a weather record"})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Weather record has been added"})
}

func (r *Repository) GetHistory(context *fiber.Ctx) error {
	userID, ok := context.Locals("user").(int)
	if !ok {
		log.Println("Failed to extract userID from context")
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid user ID"})
	}

	var weatherModels []models.Weather
	if err := r.DB.Where("user_id = ?", userID).Find(&weatherModels).Error; err != nil {
		log.Println("Error fetching search history:", err)
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not get search history"})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{"message": "History fetched successfully", "data": weatherModels})
}

func (r *Repository) DeleteWeather(context *fiber.Ctx) error {
	userID, ok := context.Locals("user").(int)
	if !ok {
		log.Println("Failed to extract userID from context")
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid user ID"})
	}

	id := context.Params("id")
	if id == "" {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "ID parameter cannot be empty"})
	}

	if err := r.DB.Where("user_id = ? AND id = ?", userID, id).Delete(&Weather{}).Error; err != nil {
		log.Printf("Error deleting weather record: %s", err)
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not delete weather", "error": err.Error()})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Weather deleted successfully"})
}

type WeatherUpdateRequest struct {
	Lat  string `json:"lat"`
	Lng  string `json:"lng"`
	Temp string `json:"temp"`
}

func (r *Repository) UpdateWeather(context *fiber.Ctx) error {
	userID, ok := context.Locals("user").(int)
	if !ok {
		log.Println("Failed to extract userID from context")
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid user ID"})
	}

	id := context.Params("id")
	if id == "" {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "ID parameter cannot be empty"})
	}

	var requestBody WeatherUpdateRequest
	if err := context.BodyParser(&requestBody); err != nil {
		log.Println("Error parsing request body:", err)
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Error parsing request body"})
	}

	updates := make(map[string]interface{})

	if requestBody.Lat != "" {
		updates["lat"], _ = strconv.ParseFloat(requestBody.Lat, 64)
	}

	if requestBody.Lng != "" {
		updates["long"], _ = strconv.ParseFloat(requestBody.Lng, 64)
	}

	if requestBody.Temp != "" {
		updates["temperature"], _ = strconv.ParseFloat(requestBody.Temp, 64)
	}

	if len(updates) == 0 {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "No valid data provided for update"})
	}

	tx := r.DB.Begin()

	if err := tx.Model(&Weather{}).Where("user_id = ?", userID).Where("id = ?", id).Updates(updates).Error; err != nil {
		tx.Rollback()
		log.Println("Error updating weather data:", err)
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Error updating weather data in database"})
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing the transaction:", err)
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Error committing the transaction"})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Weather data updated successfully"})
}
