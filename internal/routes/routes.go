package routes

import (
	"tasks-service-demo/internal/handlers"
	"tasks-service-demo/internal/middleware"
	"tasks-service-demo/internal/models"
	"tasks-service-demo/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, taskService *services.TaskService) {
	taskHandler := handlers.NewTaskHandler(taskService)

	// Health check endpoint
	app.Get("/health", handlers.HealthCheck)

	// Task API endpoints
	app.Get("tasks", taskHandler.GetAllTasks)

	app.Get("tasks/:id",
		middleware.ValidatePathID(),
		taskHandler.GetTaskByID,
	)

	app.Delete("tasks/:id",
		middleware.ValidatePathID(),
		taskHandler.DeleteTask,
	)

	app.Post("tasks",
		middleware.ValidateRequest[models.CreateTaskRequest](),
		taskHandler.CreateTask,
	)

	app.Put("tasks/:id",
		middleware.ValidatePathID(),
		middleware.ValidateRequest[models.UpdateTaskRequest](),
		taskHandler.UpdateTask,
	)
}