package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func TokenValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() != "/api/register" && c.Path() != "/api/login" {
			token := c.Get("Authorization")
			log.Println("Extracted Token: ", token)
			token = strings.TrimPrefix(token, "Bearer ")
			log.Println("Extracted Token: ", token)
			if token == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized because token is empty",
				})
			}
			claims := jwt.MapClaims{}
			jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized because jwt claims parsing failed",
					"err":     err.Error(),
				})
			}
			if jwtToken.Valid {
				return c.Next()
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized because jwtToken is InValid",
			})
		}
		return c.Next()
	}
}

// Middleware function to extract and store user ID from JWT token in the context
func ExtractUserIDFromToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() != "/api/register" && c.Path() != "/api/login" {
			// Get the JWT token from the Authorization header
			tokenString := c.Get("Authorization")

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			// Parse the JWT token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil // Provide the secret key used for signing the token
			})

			// Check for errors in token parsing
			if err != nil || !token.Valid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized at token validation",
					"err":     err.Error(),
				})
			}

			// Extract user ID from the token claims
			claims := token.Claims.(jwt.MapClaims)
			userID, ok := claims["userId"].(float64)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized at token claims",
				})
			}

			// Store the user ID in the context
			c.Locals("user", int(userID))

			// Call the next handler in the chain
			return c.Next()
		}
		return c.Next()
	}
}
