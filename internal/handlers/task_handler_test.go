package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"tasks-service-demo/internal/entities"
	"tasks-service-demo/internal/middleware"
	"tasks-service-demo/internal/requests"
	"tasks-service-demo/internal/services"
	"tasks-service-demo/internal/storage"
	"tasks-service-demo/internal/storage/naive"

	"github.com/gofiber/fiber/v2"
)

func setupTestApp() (*fiber.App, *TaskHandler) {
	app := fiber.New()
	storage.ResetStore()
	storage.InitStore(naive.NewMemoryStore())
	service := services.NewTaskService()
	handler := NewTaskHandler(service)
	return app, handler
}

func TestGetAllTasks_EmptyStore(t *testing.T) {
	app, handler := setupTestApp()
	app.Get("/tasks", handler.GetAllTasks)

	req := httptest.NewRequest("GET", "/tasks", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var tasks []entities.Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected empty tasks array, got %d tasks", len(tasks))
	}
}

func TestCreateTask_Success(t *testing.T) {
	app, handler := setupTestApp()
	app.Post("/tasks", middleware.ValidateRequest[requests.CreateTaskRequest](), handler.CreateTask)

	taskReq := requests.CreateTaskRequest{
		Name:   "Test Task",
		Status: 0,
	}
	reqBody, _ := json.Marshal(taskReq)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("Expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var task entities.Task
	if err := json.Unmarshal(body, &task); err != nil {
		t.Fatal(err)
	}

	if task.Name != "Test Task" {
		t.Errorf("Expected task name 'Test Task', got '%s'", task.Name)
	}
	if task.Status != 0 {
		t.Errorf("Expected task status 0, got %d", task.Status)
	}
	if task.ID == 0 {
		t.Error("Expected task ID to be set")
	}
}

func TestCreateTask_InvalidJSON(t *testing.T) {
	app, handler := setupTestApp()
	app.Post("/tasks", middleware.ValidateRequest[requests.CreateTaskRequest](), handler.CreateTask)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestCreateTask_ValidationError(t *testing.T) {
	app, handler := setupTestApp()
	app.Post("/tasks", middleware.ValidateRequest[requests.CreateTaskRequest](), handler.CreateTask)

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
			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
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

func TestGetTaskByID_Success(t *testing.T) {
	app, handler := setupTestApp()
	app.Get("/tasks/:id", middleware.ValidatePathID(), handler.GetTaskByID)
	app.Post("/tasks", middleware.ValidateRequest[requests.CreateTaskRequest](), handler.CreateTask)

	taskReq := requests.CreateTaskRequest{Name: "Test Task", Status: 0}
	reqBody, _ := json.Marshal(taskReq)
	createReq := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := app.Test(createReq)

	body, _ := io.ReadAll(createResp.Body)
	var createdTask entities.Task
	json.Unmarshal(body, &createdTask)

	req := httptest.NewRequest("GET", fmt.Sprintf("/tasks/%d", createdTask.ID), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	body, _ = io.ReadAll(resp.Body)
	var task entities.Task
	if err := json.Unmarshal(body, &task); err != nil {
		t.Fatal(err)
	}

	if task.ID != createdTask.ID {
		t.Errorf("Expected task ID %d, got %d", createdTask.ID, task.ID)
	}
}

func TestGetTaskByID_NotFound(t *testing.T) {
	app, handler := setupTestApp()
	app.Get("/tasks/:id", middleware.ValidatePathID(), handler.GetTaskByID)

	req := httptest.NewRequest("GET", "/tasks/999", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestGetTaskByID_InvalidID(t *testing.T) {
	app, handler := setupTestApp()
	app.Get("/tasks/:id", middleware.ValidatePathID(), handler.GetTaskByID)

	req := httptest.NewRequest("GET", "/tasks/abc", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestUpdateTask_Success(t *testing.T) {
	app, handler := setupTestApp()
	app.Post("/tasks", middleware.ValidateRequest[requests.CreateTaskRequest](), handler.CreateTask)
	app.Put("/tasks/:id", middleware.ValidatePathID(), middleware.ValidateRequest[requests.UpdateTaskRequest](), handler.UpdateTask)

	taskReq := requests.CreateTaskRequest{Name: "Original Task", Status: 0}
	reqBody, _ := json.Marshal(taskReq)
	createReq := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := app.Test(createReq)

	body, _ := io.ReadAll(createResp.Body)
	var createdTask entities.Task
	json.Unmarshal(body, &createdTask)

	updateReq := requests.UpdateTaskRequest{Name: "Updated Task", Status: 1}
	updateBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/tasks/%d", createdTask.ID), bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	body, _ = io.ReadAll(resp.Body)
	var updatedTask entities.Task
	if err := json.Unmarshal(body, &updatedTask); err != nil {
		t.Fatal(err)
	}

	if updatedTask.Name != "Updated Task" {
		t.Errorf("Expected updated name 'Updated Task', got '%s'", updatedTask.Name)
	}
	if updatedTask.Status != 1 {
		t.Errorf("Expected updated status 1, got %d", updatedTask.Status)
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	app, handler := setupTestApp()
	app.Put("/tasks/:id", middleware.ValidatePathID(), middleware.ValidateRequest[requests.UpdateTaskRequest](), handler.UpdateTask)

	updateReq := requests.UpdateTaskRequest{Name: "Updated Task", Status: 1}
	updateBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/tasks/999", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	app, handler := setupTestApp()
	app.Post("/tasks", middleware.ValidateRequest[requests.CreateTaskRequest](), handler.CreateTask)
	app.Delete("/tasks/:id", middleware.ValidatePathID(), handler.DeleteTask)

	taskReq := requests.CreateTaskRequest{Name: "Task to Delete", Status: 0}
	reqBody, _ := json.Marshal(taskReq)
	createReq := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := app.Test(createReq)

	body, _ := io.ReadAll(createResp.Body)
	var createdTask entities.Task
	json.Unmarshal(body, &createdTask)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/tasks/%d", createdTask.ID), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusNoContent {
		t.Errorf("Expected status %d, got %d", fiber.StatusNoContent, resp.StatusCode)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	app, handler := setupTestApp()
	app.Delete("/tasks/:id", middleware.ValidatePathID(), handler.DeleteTask)

	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// RESTful DELETE should be idempotent - return 204 even for non-existent resources
	if resp.StatusCode != fiber.StatusNoContent {
		t.Errorf("Expected status %d, got %d", fiber.StatusNoContent, resp.StatusCode)
	}
}
