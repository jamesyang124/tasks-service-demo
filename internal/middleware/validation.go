package middleware

import (
	"strconv"

	"tasks-service-demo/internal/models"

	"github.com/gofiber/fiber/v2"
)

func ValidateRequest[T models.Validatable]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req T

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}

		if err := req.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Locals("validated_request", req)
		return c.Next()
	}
}

func ValidatePathID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid ID",
				"message": "ID must be a valid integer",
			})
		}

		c.Locals("validated_id", id)
		return c.Next()
	}
}

func GetValidatedRequest[T models.Validatable](c *fiber.Ctx) T {
	val := c.Locals("validated_request")
	if val == nil {
		var zero T
		return zero
	}
	return val.(T)
}

func GetValidatedID(c *fiber.Ctx) int {
	val := c.Locals("validated_id")
	if val == nil {
		return 0
	}
	return val.(int)
}
