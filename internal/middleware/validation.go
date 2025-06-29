package middleware

import (
	"strconv"

	"tasks-service-demo/internal/errors"
	"tasks-service-demo/internal/requests"

	"github.com/gofiber/fiber/v2"
)

// Package middleware provides Fiber middleware for request validation and ID extraction.

// ValidateRequest returns a middleware that validates the request body against the Validatable interface.
// It parses the JSON body and validates it using the request's Validate method.
func ValidateRequest[T requests.Validatable]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req T

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&errors.ErrorResponse{
				Message: "Invalid JSON",
				Code:    errors.ErrCodeInvalidJSON,
			})
		}

		if err := req.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&errors.ErrorResponse{
				Message: err.Error(),
				Code:    errors.ErrCodeTaskInvalidInput,
			})
		}

		c.Locals("validated_request", req)
		return c.Next()
	}
}

// ValidatePathID returns a middleware that validates the :id path parameter as an integer.
// It extracts the ID from the URL path and converts it to an integer.
func ValidatePathID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&errors.ErrorResponse{
				Code:    errors.ErrCodeInvalidID,
				Message: "ID must be a valid integer",
			})
		}

		c.Locals("validated_id", id)
		return c.Next()
	}
}

// GetValidatedRequest retrieves the validated request struct from context.
// Returns the request that was previously validated by ValidateRequest middleware.
func GetValidatedRequest[T requests.Validatable](c *fiber.Ctx) T {
	val := c.Locals("validated_request")
	if val == nil {
		var zero T
		return zero
	}
	return val.(T)
}

// GetValidatedID retrieves the validated ID from context.
// Returns the ID that was previously validated by ValidatePathID middleware.
func GetValidatedID(c *fiber.Ctx) int {
	val := c.Locals("validated_id")
	if val == nil {
		return 0
	}
	return val.(int)
}
