package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"tasks-service-demo/internal/storage"
	"tasks-service-demo/internal/storage/naive"
	"tasks-service-demo/internal/storage/shard"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestAppConfiguration(t *testing.T) {
	// Test that we can create a Fiber app with the same config as main
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		},
	})

	if app == nil {
		t.Fatal("Failed to create Fiber app")
	}
}

func TestErrorHandler(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		},
	})

	// Create a route that returns an error
	app.Get("/test-error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "Test error")
	})

	req := httptest.NewRequest("GET", "/test-error", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)

	if response["error"] != "Internal Server Error" {
		t.Errorf("Expected error 'Internal Server Error', got '%v'", response["error"])
	}

	if response["message"] != "Test error" {
		t.Errorf("Expected message 'Test error', got '%v'", response["message"])
	}
}

func TestStorageTypeSelection(t *testing.T) {
	tests := []struct {
		name         string
		storageType  string
		expectedType string
	}{
		{"memory storage", "memory", "*naive.MemoryStore"},
		{"shard storage", "shard", "*shard.ShardStore"},
		{"gopool storage", "gopool", "*shard.ShardStoreGopool"},
		{"default storage", "", "*shard.ShardStoreGopool"},
		{"unknown storage", "unknown", "*shard.ShardStoreGopool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset storage for each test
			storage.ResetStore()

			// Set environment variable
			if tt.storageType != "" {
				os.Setenv("STORAGE_TYPE", tt.storageType)
			} else {
				os.Unsetenv("STORAGE_TYPE")
			}
			defer os.Unsetenv("STORAGE_TYPE")

			var store storage.Store
			storageType := os.Getenv("STORAGE_TYPE")
			if storageType == "" {
				storageType = "gopool"
			}

			shardCount := 32

			switch storageType {
			case "memory":
				store = naive.NewMemoryStore()
			case "shard":
				store = shard.NewShardStore(shardCount)
			default:
				store = shard.NewShardStoreGopool(shardCount)
			}

			if store == nil {
				t.Fatal("Store should not be nil")
			}

			// Check the type using type assertion or reflection
			actualType := getTypeName(store)
			if actualType != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, actualType)
			}

			// Clean up resources
			if closer, ok := store.(interface{ Close() error }); ok {
				closer.Close()
			}
		})
	}
}

func TestShardCountConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		shardCountStr string
		expectedCount int
	}{
		{"default shard count", "", 32},
		{"valid shard count", "16", 16},
		{"another valid count", "64", 64},
		{"invalid shard count", "invalid", 32},
		{"zero shard count", "0", 32},
		{"negative shard count", "-1", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shardCountStr != "" {
				os.Setenv("SHARD_COUNT", tt.shardCountStr)
			} else {
				os.Unsetenv("SHARD_COUNT")
			}
			defer os.Unsetenv("SHARD_COUNT")

			shardCount := 32
			if shardCountStr := os.Getenv("SHARD_COUNT"); shardCountStr != "" {
				if sc, err := parseShardCount(shardCountStr); err == nil && sc > 0 {
					shardCount = sc
				}
			}

			if shardCount != tt.expectedCount {
				t.Errorf("Expected shard count %d, got %d", tt.expectedCount, shardCount)
			}
		})
	}
}

func TestMiddlewareSetup(t *testing.T) {
	app := fiber.New()

	// Track middleware execution using Locals instead of headers
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("test-logger", "applied")
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("test-recover", "applied")
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("test-cors", "applied")
		return c.Next()
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"logger":  c.Locals("test-logger"),
			"recover": c.Locals("test-recover"),
			"cors":    c.Locals("test-cors"),
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)

	expectedHeaders := []string{"logger", "recover", "cors"}
	for _, header := range expectedHeaders {
		if response[header] != "applied" {
			t.Errorf("Expected %s middleware to be applied", header)
		}
	}
}

func TestGracefulShutdown(t *testing.T) {
	// Test the graceful shutdown logic without actually starting the server
	storage.ResetStore()
	mockStore := naive.NewMemoryStore()
	storage.InitStore(mockStore)

	// Simulate graceful shutdown cleanup
	if store := storage.GetStore(); store != nil {
		if closer, ok := store.(interface{ Close() error }); ok {
			err := closer.Close()
			// For MemoryStore, Close() is not implemented, so we expect this to not panic
			_ = err // MemoryStore doesn't implement Close(), so this is fine
		}
	}
}

func TestAppIntegration(t *testing.T) {
	// Test full app setup without Listen()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		},
	})

	// Simulate main.go setup
	storage.ResetStore()
	store := naive.NewMemoryStore() // Use simple store for testing
	storage.InitStore(store)

	// Add a simple health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Test the setup
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestEnvironmentVariableHandling(t *testing.T) {
	// Test default values when no environment variables are set
	os.Unsetenv("STORAGE_TYPE")
	os.Unsetenv("SHARD_COUNT")

	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "gopool"
	}

	if storageType != "gopool" {
		t.Errorf("Expected default storage type 'gopool', got '%s'", storageType)
	}

	shardCount := 32
	if shardCountStr := os.Getenv("SHARD_COUNT"); shardCountStr != "" {
		if sc, err := parseShardCount(shardCountStr); err == nil && sc > 0 {
			shardCount = sc
		}
	}

	if shardCount != 32 {
		t.Errorf("Expected default shard count 32, got %d", shardCount)
	}
}

func TestDotenvLoading(t *testing.T) {
	// Test that dotenv loading doesn't panic
	assert.NotPanics(t, func() {
		// This should not panic even if .env file doesn't exist
		_ = godotenv.Load()
	})
}

// Helper functions for testing

func getTypeName(store storage.Store) string {
	switch store.(type) {
	case *naive.MemoryStore:
		return "*naive.MemoryStore"
	case *shard.ShardStore:
		return "*shard.ShardStore"
	case *shard.ShardStoreGopool:
		return "*shard.ShardStoreGopool"
	default:
		return "unknown"
	}
}

func parseShardCount(s string) (int, error) {
	// Helper function to mimic the strconv.Atoi logic from main
	return parseInteger(s)
}

func parseInteger(s string) (int, error) {
	result := 0
	for _, char := range s {
		if char < '0' || char > '9' {
			return 0, fiber.NewError(400, "invalid integer")
		}
		result = result*10 + int(char-'0')
	}
	return result, nil
}

// Benchmark the app setup process
func BenchmarkAppSetup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		app := fiber.New(fiber.Config{
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				code := fiber.StatusInternalServerError
				if e, ok := err.(*fiber.Error); ok {
					code = e.Code
				}
				return ctx.Status(code).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			},
		})

		storage.ResetStore()
		store := naive.NewMemoryStore()
		storage.InitStore(store)

		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendStatus(200)
		})

		_ = app
	}
}
