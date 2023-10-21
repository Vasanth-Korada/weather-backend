package main

import (
	"os"

	"github.com/Vasanth-Korada/weather-backend/storage"
	"gorm.io/gorm"

	"github.com/Vasanth-Korada/weather-backend/models"
	"github.com/Vasanth-Korada/weather-backend/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var secretKey = []byte(os.Getenv("gloresoft-samson"))

func initializeConfig() *storage.Config {
	return &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

func handleDBMigrations(db *gorm.DB) {
	errWeather := models.MigrateWeather(db)
	HandleError(errWeather, "Could not migrate database")

	errUsers := models.MigrateUsers(db)
	HandleError(errUsers, "Could not migrate database")

}

func main() {
	err := godotenv.Load(".env")
	HandleError(err, "Error loading .env file")

	config := initializeConfig()

	db, err := storage.NewConnection(config)
	HandleError(err, "Could not load the database")

	handleDBMigrations(db)

	r := repository.Repository{
		DB: db,
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Use(TokenValidationMiddleware())
	app.Use(ExtractUserIDFromToken())

	r.SetupRoutes(app)
	app.Listen(":8080")
}
