package handlers

import (
	apperrors "tasks-service-demo/internal/errors"
	"tasks-service-demo/internal/middleware"
	"tasks-service-demo/internal/requests"
	"tasks-service-demo/internal/services"

	"github.com/gofiber/fiber/v2"
)

// Package handlers provides HTTP handlers for the Task API.

// TaskHandler handles HTTP requests for task operations.
type TaskHandler struct {
	service *services.TaskService
}

// NewTaskHandler creates a new TaskHandler with the given TaskService.
func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// GetAllTasks handles GET /tasks and returns all tasks.
func (h *TaskHandler) GetAllTasks(c *fiber.Ctx) error {
	tasks := h.service.GetAllTasks()
	return c.JSON(tasks)
}

// GetTaskByID handles GET /tasks/:id and returns a task by its ID.
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id := middleware.GetValidatedID(c)

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		switch err.Code {
		case apperrors.ErrCodeTaskNotFound:
			return c.Status(fiber.StatusBadRequest).JSON(apperrors.ToResponse(err))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(apperrors.ErrInternalErrorResponse)
		}
	}

	return c.JSON(task)
}

// CreateTask handles POST /tasks and creates a new task.
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	req := middleware.GetValidatedRequest[requests.CreateTaskRequest](c)

	task, err := h.service.CreateTask(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(apperrors.ErrInternalErrorResponse)
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// UpdateTask handles PUT /tasks/:id and updates an existing task.
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := middleware.GetValidatedID(c)
	req := middleware.GetValidatedRequest[requests.UpdateTaskRequest](c)

	task, err := h.service.UpdateTask(id, &req)
	if err != nil {
		switch err.Code {
		case apperrors.ErrCodeTaskNotFound:
			return c.Status(fiber.StatusBadRequest).JSON(apperrors.ToResponse(err))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(apperrors.ErrInternalErrorResponse)
		}
	}

	return c.JSON(task)
}

// DeleteTask handles DELETE /tasks/:id and deletes a task by its ID.
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := middleware.GetValidatedID(c)

	err := h.service.DeleteTask(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(apperrors.ErrInternalErrorResponse)
	}

	// RESTful DELETE: Always return 204 No Content for successful DELETE (idempotent)
	return c.Status(fiber.StatusNoContent).Send(nil)
}
