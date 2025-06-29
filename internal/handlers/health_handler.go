package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// Package handlers provides HTTP handlers for the Task API.

// HealthCheck handles GET /health and returns the health status of the API.
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Task API is running",
	})
}
