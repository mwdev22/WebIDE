package handlers

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type FiberHandler func(*fiber.Ctx) error

func AuthMiddleware(handler FiberHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return Unauthorized("missing authorization header")
		}

		if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
			return Unauthorized("invalid authorization header format")
		}

		tokenStr := authHeader[7:]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// verify the token's signature method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			fmt.Println("Error:", err)
			return Unauthorized("invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return Unauthorized("invalid token claims")
		}

		userID, ok := claims["userID"].(float64)
		if !ok {
			return Unauthorized("error parsing userID from claims")
		}

		c.Locals("userID", uint(userID))

		return handler(c)
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
