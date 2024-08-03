package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ApiError struct {
	StatusCode int `json:"status_code"`
	Msg        any `json:"msg"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("api error: %v", e.Msg)
}

// base api error
func NewApiError(statusCode int, e error) ApiError {
	return ApiError{
		StatusCode: statusCode,
		Msg:        e.Error(),
	}
}

func InvalidJSON() ApiError {
	return NewApiError(fiber.StatusUnprocessableEntity, fmt.Errorf("invalid json data"))
}

func BadQuery(err error) ApiError {
	return NewApiError(fiber.StatusInternalServerError, err)
}

func ExternalServiceErr(err error) ApiError {
	return NewApiError(fiber.StatusBadRequest, err)
}

func Unauthorized(msg any) ApiError {
	return NewApiError(fiber.StatusUnauthorized, fmt.Errorf("Unauthorized, %s", msg))
}
