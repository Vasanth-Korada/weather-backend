package repository

import (
	"errors"
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

var secretKey []byte

func generateToken(userId int) (string, error) {
	secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(secretKey) == 0 {
		return "", errors.New("secret key is empty")
	}

	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (r *Repository) handleLogin(c *fiber.Ctx) error {
	var creds Credentials
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	user, err := models.AuthenticateUser(r.DB, creds.Username, creds.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid username or password"})
	}

	token, err := generateToken(user.UserID)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate token"})
	}

	return c.SendString(token)
}

func (r *Repository) handleRegister(c *fiber.Ctx) error {
	var creds Credentials
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if creds.Username == "" || creds.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Username and Password are required"})
	}

	err := models.RegisterUser(r.DB, creds.Username, creds.Password, creds.DOB)
	if err != nil {
		log.Printf("User registration error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to register user"})
	}

	return c.Status(fiber.StatusCreated).SendString("User registered successfully")
}

func (r *Repository) handleLogout(c *fiber.Ctx) error {
	c.Response().Header.Del("Authorization")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}
