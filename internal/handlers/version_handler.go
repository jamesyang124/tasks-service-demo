package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

// Package handlers provides HTTP handlers for the Task API.

// VersionInfo represents the version information of the API
type VersionInfo struct {
	Version string `json:"version"`
}

// VersionHandler returns version information about the API
func VersionHandler(c *fiber.Ctx) error {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "1.0.0" // Default semantic version
	}

	versionInfo := VersionInfo{
		Version: version,
	}

	return c.JSON(versionInfo)
}
