package main

import (
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	apperrors "tasks-service-demo/internal/errors"
	applog "tasks-service-demo/internal/logger"
	"tasks-service-demo/internal/routes"
	"tasks-service-demo/internal/services"
	"tasks-service-demo/internal/storage"
	"tasks-service-demo/internal/storage/naive"
	"tasks-service-demo/internal/storage/shard"
	"tasks-service-demo/internal/storage/xsync"
)

func main() {
	// flush zap on exit
	defer func() {
		if err := applog.Get().Sync(); err != nil {
			applog.Get().Warnf("Zap sync/flush error: %v", err)
		}
	}()

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		applog.Get().Info("No .env file found, using system environment variables")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"code":    apperrors.ErrCodeInternalError,
				"message": err.Error(),
			})
		},
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Initialize storage with configuration options
	var store storage.Store

	// Check environment variable for storage type (default: xsync)
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "xsync" // Default to lock-free best performance
	}

	// Configure shard count (default: 32 for M4 Pro optimization)
	shardCount := 32
	if shardCountStr := os.Getenv("SHARD_COUNT"); shardCountStr != "" {
		if sc, err := strconv.Atoi(shardCountStr); err == nil && sc > 0 {
			shardCount = sc
		}
	}

	// Initialize based on configuration
	switch storageType {
	case "xsync":
		store = xsync.NewXSyncStore()
		applog.Get().Info("XSyncStore initialized (lock-free concurrent map - best performance)")
	case "gopool":
		store = shard.NewShardStoreGopool(shardCount)
		applog.Get().Infof("ShardStoreGopool initialized with %d shards", shardCount)
	case "shard":
		store = shard.NewShardStore(shardCount)
		applog.Get().Infof("ShardStore initialized with dedicated workers and %d shards", shardCount)
	case "memory":
		store = naive.NewMemoryStore()
		applog.Get().Info("MemoryStore initialized (single mutex - not recommended for production)")
	default:
		// Default to xsync for best performance
		store = xsync.NewXSyncStore()
		applog.Get().Infof("Unknown storage type '%s', defaulting to XSyncStore", storageType)
		applog.Get().Info("XSyncStore initialized (lock-free concurrent map - best performance)")
	}

	storage.InitStore(store)
	taskService := services.NewTaskService()
	routes.SetupRoutes(app, taskService)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	// Graceful shutdown with proper resource cleanup
	go func() {
		defer wg.Done()

		<-quit
		applog.Get().Info("Received shutdown signal...")

		// Gracefully shutdown Fiber after 5 seconds
		if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
			applog.Get().Errorf("Fiber shutdown error: %v", err)
		} else {
			applog.Get().Info("Fiber server shutdown complete")
		}

		// Close storage resources before shutting down server
		if store := storage.GetStore(); store != nil {
			if closer, ok := store.(interface{ Close() error }); ok {
				if err := closer.Close(); err != nil {
					applog.Get().Errorf("Error closing storage: %v", err)
				} else {
					applog.Get().Info("Storage resources cleaned up")
				}
			}
		}

	}()

	applog.Get().Info("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		applog.Get().Fatalf("Server failed to start: %v", err)
	}

	// Wait for cleanup goroutine to complete
	wg.Wait()
	applog.Get().Info("Server gracefully stopped")
}
