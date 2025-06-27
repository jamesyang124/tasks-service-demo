package handlers

import (
	"strings"
	"tasks-service-demo/internal/middleware"
	"tasks-service-demo/internal/models"
	"tasks-service-demo/internal/services"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetAllTasks(c *fiber.Ctx) error {
	tasks := h.service.GetAllTasks()
	return c.JSON(tasks)
}

func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id := middleware.GetValidatedID(c)

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(&models.ErrorResponse{
				Error:   "Task not found",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&models.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
	}

	return c.JSON(task)
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	req := middleware.GetValidatedRequest[models.CreateTaskRequest](c)

	task, err := h.service.CreateTask(&req)
	if err != nil {
		if _, ok := err.(*models.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(&models.ErrorResponse{
				Error:   "Validation error",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&models.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := middleware.GetValidatedID(c)
	req := middleware.GetValidatedRequest[models.UpdateTaskRequest](c)

	task, err := h.service.UpdateTask(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(&models.ErrorResponse{
				Error:   "Task not found",
				Message: err.Error(),
			})
		}
		if _, ok := err.(*models.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(&models.ErrorResponse{
				Error:   "Validation error",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&models.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
	}

	return c.JSON(task)
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := middleware.GetValidatedID(c)

	err := h.service.DeleteTask(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(&models.ErrorResponse{
				Error:   "Task not found",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&models.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}