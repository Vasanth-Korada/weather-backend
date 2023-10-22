package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func extractToken(c *fiber.Ctx) string {
	token := c.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	return token
}

func TokenValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() != "/api/register" && c.Path() != "/api/login" {
			tokenString := extractToken(c)
			if tokenString == "" {
				log.Println("Token is empty")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized",
				})
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.NewValidationError("Invalid signing method", jwt.ValidationErrorSignatureInvalid)
				}
				return secretKey, nil
			})

			if err != nil || !token.Valid {
				log.Printf("Failed to validate the token: %v", err)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized",
				})
			}
		}
		return c.Next()
	}
}

func ExtractUserIDFromToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() != "/api/register" && c.Path() != "/api/login" {
			tokenString := extractToken(c)
			token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("Failed to extract claims from the token")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized",
				})
			}

			userID, ok := claims["userId"].(float64)
			if !ok {
				log.Println("UserID not found in the token claims")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized",
				})
			}

			c.Locals("user", int(userID))
		}
		return c.Next()
	}
}
