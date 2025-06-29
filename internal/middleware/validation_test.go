package middleware

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"tasks-service-demo/internal/requests"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupTestApp() *fiber.App {
	return fiber.New()
}

func TestValidateRequest_Success(t *testing.T) {
	app := setupTestApp()

	app.Post("/test", ValidateRequest[requests.CreateTaskRequest](), func(c *fiber.Ctx) error {
		req := GetValidatedRequest[requests.CreateTaskRequest](c)
		return c.JSON(fiber.Map{
			"received": req,
		})
	})

	taskReq := requests.CreateTaskRequest{
		Name:   "Test Task",
		Status: 0,
	}
	reqBody, _ := json.Marshal(taskReq)

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestValidateRequest_InvalidJSON(t *testing.T) {
	app := setupTestApp()

	app.Post("/test", ValidateRequest[requests.CreateTaskRequest](), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestValidateRequest_ValidationError(t *testing.T) {
	app := setupTestApp()

	app.Post("/test", ValidateRequest[requests.CreateTaskRequest](), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})

	tests := []struct {
		name string
		req  requests.CreateTaskRequest
	}{
		{"empty name", requests.CreateTaskRequest{Name: "", Status: 0}},
		{"invalid status", requests.CreateTaskRequest{Name: "Test", Status: 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.req)
			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != fiber.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
			}
		})
	}
}

func TestValidatePathID_Success(t *testing.T) {
	app := setupTestApp()

	app.Get("/test/:id", ValidatePathID(), func(c *fiber.Ctx) error {
		id := GetValidatedID(c)
		return c.JSON(fiber.Map{
			"id": id,
		})
	})

	req := httptest.NewRequest("GET", "/test/123", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestValidatePathID_InvalidID(t *testing.T) {
	app := setupTestApp()

	app.Get("/test/:id", ValidatePathID(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})

	tests := []string{
		"/test/abc",
		"/test/12.34",
		"/test/notanumber",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != fiber.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
			}
		})
	}
}

func TestGetValidatedRequest_WithoutMiddleware(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		req := GetValidatedRequest[requests.CreateTaskRequest](c)
		// Should return zero value when no validated request is stored
		if req.Name != "" || req.Status != 0 {
			t.Error("Expected zero value for request when middleware not used")
		}
		return c.JSON(fiber.Map{"success": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestGetValidatedID_WithoutMiddleware(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		id := GetValidatedID(c)
		// Should return 0 when no validated ID is stored
		if id != 0 {
			t.Errorf("Expected 0 for ID when middleware not used, got %d", id)
		}
		return c.JSON(fiber.Map{"success": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}
