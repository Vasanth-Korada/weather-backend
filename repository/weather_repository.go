package repository

import (
	"fmt"
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
	userID := context.Locals("user").(int)
	fmt.Println(userID)
	weather := Weather{}

	err := context.BodyParser(&weather)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}
	weather.UserID = userID

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
	userID := context.Locals("user").(int)
	weatherModels := &[]models.Weather{}
	err := r.DB.Where("user_id = ?", userID).Find(weatherModels).Error

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
	userID := context.Locals("user").(int)
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	if err := r.DB.Where("user_id = ? AND id = ?", userID, id).Delete(&Weather{}).Error; err != nil {
		log.Printf("Error deleting weather record: %s", err.Error())
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "could not delete weather",
			"error":   err.Error(),
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "weather deleted successfully",
	})
	return nil
}

type WeatherUpdateRequest struct {
	Lat  string `json:"lat"`
	Lng  string `json:"lng"`
	Temp string `json:"temp"`
}

func (r *Repository) UpdateWeather(context *fiber.Ctx) error {
	userID := context.Locals("user").(int)
	id := context.Params("id")

	if id == "" {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "ID parameter cannot be empty",
		})
	}

	var requestBody WeatherUpdateRequest
	if err := context.BodyParser(&requestBody); err != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Error parsing request body",
		})
	}

	updates := make(map[string]interface{})

	if requestBody.Lat != "" { // using 0 as float64 default value
		updates["lat"], _ = strconv.ParseFloat(requestBody.Lat, 64)
	}

	if requestBody.Lng != "" {
		updates["long"], _ = strconv.ParseFloat(requestBody.Lng, 64)
	}

	if requestBody.Temp != "" {
		updates["temperature"], _ = strconv.ParseFloat(requestBody.Temp, 64)
	}

	// Check if there's any data to update
	if len(updates) == 0 {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "No valid data provided for update",
		})
	}

	// Starting a transaction
	tx := r.DB.Begin()

	err := tx.Model(&Weather{}).Where("user_id = ?", userID).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		tx.Rollback()
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error updating weather data in database",
		})
	}

	if err := tx.Commit().Error; err != nil {
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error committing the transaction",
		})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Weather data updated successfully",
	})
}
