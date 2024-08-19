package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mwdev22/WebIDE/cmd/types"
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

func SQLError(err error) ApiError {
	return NewApiError(fiber.StatusInternalServerError, err)
}

func ExternalServiceErr(err error) ApiError {
	return NewApiError(fiber.StatusBadRequest, err)
}

func Unauthorized(msg any) ApiError {
	return NewApiError(fiber.StatusUnauthorized, fmt.Errorf("unauthorized, %s", msg))
}

func BadQueryParameter(name string) ApiError {
	return NewApiError(fiber.StatusBadRequest, fmt.Errorf("bad query param, %s", name))
}

func ValidationError(errors []*types.ErrorResponse) ApiError {
	return ApiError{
		StatusCode: fiber.StatusBadRequest,
		Msg:        errors,
	}
}

func NotFound(id int, name string) ApiError {
	return NewApiError(fiber.StatusNotFound, fmt.Errorf("%s with id %v not found in database", name, id))
}
