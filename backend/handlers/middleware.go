package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/mwdev22/WebIDE/backend/utils"
)

var jwtSecret = []byte(utils.SecretKey)

type FiberHandler func(*fiber.Ctx) error

func AuthMiddleware(handler FiberHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		handler(c)
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return Unauthorized("missing authorization header")
		}

		tokenStr := authHeader[len("Bearer "):]
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return Unauthorized("invalid token")
		}

		username, ok := (*claims)["username"].(string)
		if !ok {
			return Unauthorized("invalid token claims")
		}

		c.Locals("username", username)
		return c.Next()
	}
}

func ErrMiddleware(handler FiberHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := handler(c); err != nil {
			// embed returning specific errors in request
			if apiErr, ok := err.(ApiError); ok {
				return c.Status(apiErr.StatusCode).JSON(fiber.Map{
					"msg": apiErr.Msg,
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"msg": err.Error(),
				})
			}
		}
		return nil
	}
}
