package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Vasanth-Korada/weather-backend/storage"
	"gorm.io/gorm"

	"github.com/Vasanth-Korada/weather-backend/models"
	"github.com/Vasanth-Korada/weather-backend/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var secretKey []byte

func initializeConfig() (*storage.Config, error) {
	requiredEnv := []string{"DB_HOST", "DB_PORT", "DB_PASS", "DB_USER", "DB_SSLMODE", "DB_NAME", "JWT_SECRET_KEY"}

	for _, v := range requiredEnv {
		if os.Getenv(v) == "" {
			return nil, fmt.Errorf("environment variable %s is not set", v)
		}
	}

	secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	return &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}, nil
}

func handleDBMigrations(db *gorm.DB) error {
	if err := models.MigrateWeather(db); err != nil {
		return fmt.Errorf("could not migrate Weather table: %v", err)
	}

	if err := models.MigrateUsers(db); err != nil {
		return fmt.Errorf("could not migrate Users table: %v", err)
	}

	return nil
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config, err := initializeConfig()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatalf("Could not establish database connection: %v", err)
	}

	if err := handleDBMigrations(db); err != nil {
		log.Fatalf("Database migration error: %v", err)
	}

	r := repository.Repository{DB: db}

	app := fiber.New()
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))
	app.Use(TokenValidationMiddleware())
	app.Use(ExtractUserIDFromToken())

	r.SetupRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
