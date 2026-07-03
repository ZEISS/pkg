package fiberx

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

type httpError struct {
	StatusCode int
	Message    string
}

// Errors returns a Fiber error handler that formats errors into a consistent JSON structure.
func Errors() fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		var targetErr *fiber.Error
		if errors.As(err, &targetErr) {
			code = targetErr.Code
		}

		return c.Status(code).JSON(&httpError{
			StatusCode: code,
			Message:    err.Error(),
		})
	}
}
