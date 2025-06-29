package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestVersionHandler_DefaultVersion(t *testing.T) {
	// Ensure APP_VERSION env var is not set
	os.Unsetenv("APP_VERSION")

	app := fiber.New()
	app.Get("/version", VersionHandler)

	req := httptest.NewRequest("GET", "/version", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var versionInfo VersionInfo
	err = json.Unmarshal(body, &versionInfo)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if versionInfo.Version != "1.0.0" {
		t.Errorf("Expected default version '1.0.0', got '%s'", versionInfo.Version)
	}
}

func TestVersionHandler_CustomVersion(t *testing.T) {
	// Set custom version
	os.Setenv("APP_VERSION", "2.3.1")
	defer os.Unsetenv("APP_VERSION")

	app := fiber.New()
	app.Get("/version", VersionHandler)

	req := httptest.NewRequest("GET", "/version", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var versionInfo VersionInfo
	err = json.Unmarshal(body, &versionInfo)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if versionInfo.Version != "2.3.1" {
		t.Errorf("Expected version '2.3.1', got '%s'", versionInfo.Version)
	}
}

func TestVersionHandler_ResponseFormat(t *testing.T) {
	app := fiber.New()
	app.Get("/version", VersionHandler)

	req := httptest.NewRequest("GET", "/version", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// Verify content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Verify response structure
	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	// Check that version field exists
	if _, exists := response["version"]; !exists {
		t.Error("Required field 'version' missing from response")
	}
}
