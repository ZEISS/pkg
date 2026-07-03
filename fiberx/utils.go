package fiberx

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

// ParseBody is helper function for parsing the body.
// Is any error occurs it will panic.
// Its just a helper function to avoid writing if condition again n again.
func ParseBody(ctx fiber.Ctx, body interface{}) *fiber.Error {
	if err := ctx.Bind().Body(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}

// ParseAll is helper function for parsing the query and body.
func ParseAll(ctx fiber.Ctx, out interface{}) *fiber.Error {
	if err := ctx.Bind().All(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}

// ParseBodyAndValidate is helper function for parsing and validating the body.
func ParseBodyAndValidate(ctx fiber.Ctx, body interface{}) *fiber.Error {
	if err := ctx.Bind().Body(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return Validate(body)
}

// ParseAllAndValidate is helper function for parsing and validating the query and body.
func ParseAllAndValidate(ctx fiber.Ctx, out interface{}) *fiber.Error {
	if err := ctx.Bind().All(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return Validate(out)
}

// Validate validates the input struct.
func Validate(payload interface{}) *fiber.Error {
	err := validate.Struct(payload)
	if err == nil {
		return nil
	}

	var targetErr *fiber.Error
	if errors.As(err, &targetErr) {
		return fiber.NewError(fiber.StatusBadRequest, targetErr.Error())
	}

	return nil
}
