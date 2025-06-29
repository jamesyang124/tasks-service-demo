package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"tasks-service-demo/internal/routes"
	"tasks-service-demo/internal/services"
	"tasks-service-demo/internal/storage"
	"tasks-service-demo/internal/storage/bigcache"
	"tasks-service-demo/internal/storage/naive"
	"tasks-service-demo/internal/storage/shard"
)

func main() {
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

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Initialize storage with configuration options
	var store storage.Store

	// Check environment variable for storage type (default: gopool)
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "gopool" // Default to best performance
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
	case "memory":
		store = naive.NewMemoryStore()
		log.Printf("MemoryStore initialized (single mutex - not recommended for production)")
	case "bigcache":
		store = bigcache.NewBigCacheStore()
		log.Printf("BigCacheStore initialized (off-heap cache with zero GC overhead)")
	case "shard":
		store = shard.NewShardStore(shardCount)
		log.Printf("ShardStore initialized with dedicated workers and %d shards", shardCount)

	default:
		// Default to gopool for best performance
		store = shard.NewShardStoreGopool(shardCount)
		log.Printf("Unknown storage type '%s', defaulting to ShardStoreGopool", storageType)
		log.Printf("Optimized for M4 Pro 14-core architecture with %d shards", shardCount)
	}

	storage.InitStore(store)
	taskService := services.NewTaskService()
	routes.SetupRoutes(app, taskService)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Graceful shutdown with proper resource cleanup
	go func() {
		<-c
		log.Println("Gracefully shutting down...")

		// Close storage resources before shutting down server
		if store := storage.GetStore(); store != nil {
			if closer, ok := store.(interface{ Close() error }); ok {
				if err := closer.Close(); err != nil {
					log.Printf("Error closing storage: %v", err)
				} else {
					log.Println("Storage resources cleaned up")
				}
			}
		}

		_ = app.Shutdown()
	}()

	log.Println("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Server gracefully stopped")
}
