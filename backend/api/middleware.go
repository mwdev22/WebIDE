package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/mwdev22/WebIDE/backend/utils"
)

var jwtSecret = []byte(utils.SecretKey)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(map[string]string{"error": "missing authorization header"})
		}

		tokenStr := authHeader[len("Bearer "):]
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(map[string]string{"error": "invalid token"})
		}

		username, ok := (*claims)["username"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(map[string]string{"error": "invalid token claims"})
		}

		c.Locals("username", username)
		return c.Next()
	}
}
