package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"tasks-service-demo/internal/routes"
	"tasks-service-demo/internal/services"
	"tasks-service-demo/internal/storage"
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

	// Initialize storage singleton with shard store (8 shards for better distribution)
	shardStore := storage.NewShardStore(8)
	storage.InitStore(shardStore)
	taskService := services.NewTaskService()
	routes.SetupRoutes(app, taskService)
	
	// Log shard store initialization
	if ss, ok := storage.GetStore().(*storage.ShardStore); ok {
		stats := ss.GetShardStats()
		log.Printf("Shard Store initialized with %d shards", stats["numShards"])
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// this may still close early than main goroutine
	// main should have some approach wait for this down
	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	log.Println("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Server gracefully stopped")
}
