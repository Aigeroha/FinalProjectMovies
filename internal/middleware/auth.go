package middleware

import (
	"strings"
	"github.com/gofiber/fiber/v3"
	"final-project/internal/auth"
	"final-project/internal/responses"
)


func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return responses.Error(c, 401, "missing authorization token")
		}

		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return responses.Error(c, 401, "invalid authorization format")
		}

		tokenString := parts[1]

		
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			return responses.Error(c, 401, "unauthorized: token is expired or invalid")
		}

		
		c.Locals("customer_id", claims.CustomerID)

		
		return c.Next()
	}
}