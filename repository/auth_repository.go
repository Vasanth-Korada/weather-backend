package repository

import (
	"log"
	"os"
	"time"

	"github.com/Vasanth-Korada/weather-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DOB      string `json:"dob"`
}

var secretKey = []byte(os.Getenv("gloresoft-samson"))

func generateToken(userId int) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (r *Repository) handleLogin(c *fiber.Ctx) error {
	var creds Credentials
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	user, err := models.AuthenticateUser(r.DB, creds.Username, creds.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	token, err := generateToken(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.SendString(token)
}

func (r *Repository) handleRegister(c *fiber.Ctx) error {
	log.Println("Request Headers:", c.Request().Header)
	var creds Credentials
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	err := models.RegisterUser(r.DB, creds.Username, creds.Password, creds.DOB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}
	log.Println("Request Body:", creds)

	return c.Status(fiber.StatusCreated).SendString("User registered successfully")
}

func (r *Repository) handleLogout(c *fiber.Ctx) error {
	c.Response().Header.Del("Authorization")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}
