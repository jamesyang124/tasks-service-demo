package routes

import (
	"tasks-service-demo/internal/handlers"
	"tasks-service-demo/internal/middleware"
	"tasks-service-demo/internal/requests"
	"tasks-service-demo/internal/services"

	"github.com/gofiber/fiber/v2"
)

// Package routes defines the application's HTTP route setup.

// SetupRoutes registers all API routes and handlers with the Fiber app.
func SetupRoutes(app *fiber.App, taskService *services.TaskService) {
	taskHandler := handlers.NewTaskHandler(taskService)

	// Health check endpoint
	app.Get("/health", handlers.HealthCheck)

	// Version endpoint
	app.Get("/version", handlers.VersionHandler)

	// Task API endpoints
	app.Get("/tasks", taskHandler.GetAllTasks)

	app.Get("/tasks/:id",
		middleware.ValidatePathID(),
		taskHandler.GetTaskByID,
	)

	app.Delete("/tasks/:id",
		middleware.ValidatePathID(),
		taskHandler.DeleteTask,
	)

	app.Post("/tasks",
		middleware.ValidateRequest[requests.CreateTaskRequest](),
		taskHandler.CreateTask,
	)

	app.Put("/tasks/:id",
		middleware.ValidatePathID(),
		middleware.ValidateRequest[requests.UpdateTaskRequest](),
		taskHandler.UpdateTask,
	)
}
