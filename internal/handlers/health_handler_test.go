package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	app.Get("/health", HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedStatus := "ok"
	if response["status"] != expectedStatus {
		t.Errorf("Expected status '%s', got '%v'", expectedStatus, response["status"])
	}

	expectedMessage := "Task API is running"
	if response["message"] != expectedMessage {
		t.Errorf("Expected message '%s', got '%v'", expectedMessage, response["message"])
	}
}

func TestHealthCheck_ResponseFormat(t *testing.T) {
	app := fiber.New()
	app.Get("/health", HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)
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

	// Check that required fields exist
	requiredFields := []string{"status", "message"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Required field '%s' missing from response", field)
		}
	}
}
