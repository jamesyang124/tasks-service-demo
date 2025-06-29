package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"tasks-service-demo/internal/entities"
	"tasks-service-demo/internal/errors"
	"tasks-service-demo/internal/requests"
	"tasks-service-demo/internal/services"
	"tasks-service-demo/internal/storage"
	"tasks-service-demo/internal/storage/naive"

	"github.com/gofiber/fiber/v2"
)

func setupTestApp() *fiber.App {
	// Setup storage
	storage.ResetStore()
	storage.InitStore(naive.NewMemoryStore())

	// Create app and service
	app := fiber.New()
	taskService := services.NewTaskService()

	// Setup routes
	SetupRoutes(app, taskService)

	return app
}

func TestSetupRoutes_HealthEndpoint(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/health", nil)
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

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%v'", response["status"])
	}
}

func TestSetupRoutes_GetAllTasks(t *testing.T) {
	app := setupTestApp()

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
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Should return empty array initially
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

func TestSetupRoutes_CreateTask(t *testing.T) {
	app := setupTestApp()

	taskReq := requests.CreateTaskRequest{
		Name:   "Test Route Task",
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
	err = json.Unmarshal(body, &task)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if task.Name != taskReq.Name {
		t.Errorf("Expected name '%s', got '%s'", taskReq.Name, task.Name)
	}
	if task.ID == 0 {
		t.Error("Expected task ID to be set")
	}
}

func TestSetupRoutes_CreateTask_ValidationError(t *testing.T) {
	app := setupTestApp()

	// Invalid request - empty name
	taskReq := requests.CreateTaskRequest{
		Name:   "",
		Status: 0,
	}
	reqBody, _ := json.Marshal(taskReq)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errorResp errors.ErrorResponse
	err = json.Unmarshal(body, &errorResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errorResp.Code != errors.ErrCodeTaskInvalidInput {
		t.Errorf("Expected error code %d, got %d", errors.ErrCodeTaskInvalidInput, errorResp.Code)
	}
}

func TestSetupRoutes_GetTaskByID(t *testing.T) {
	app := setupTestApp()

	// First create a task
	taskReq := requests.CreateTaskRequest{Name: "Get Test Task", Status: 0}
	reqBody, _ := json.Marshal(taskReq)
	createReq := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := app.Test(createReq)

	body, _ := io.ReadAll(createResp.Body)
	var createdTask entities.Task
	json.Unmarshal(body, &createdTask)

	// Now get the task by ID
	req := httptest.NewRequest("GET", fmt.Sprintf("/tasks/%d", createdTask.ID), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	body, _ = io.ReadAll(resp.Body)
	var retrievedTask entities.Task
	err = json.Unmarshal(body, &retrievedTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if retrievedTask.ID != createdTask.ID {
		t.Errorf("Expected ID %d, got %d", createdTask.ID, retrievedTask.ID)
	}
}

func TestSetupRoutes_GetTaskByID_InvalidID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/tasks/abc", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errorResp errors.ErrorResponse
	err = json.Unmarshal(body, &errorResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errorResp.Code != errors.ErrCodeInvalidID {
		t.Errorf("Expected error code %d, got %d", errors.ErrCodeInvalidID, errorResp.Code)
	}
}

func TestSetupRoutes_GetTaskByID_NotFound(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/tasks/999", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestSetupRoutes_UpdateTask(t *testing.T) {
	app := setupTestApp()

	// First create a task
	taskReq := requests.CreateTaskRequest{Name: "Update Test Task", Status: 0}
	reqBody, _ := json.Marshal(taskReq)
	createReq := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := app.Test(createReq)

	body, _ := io.ReadAll(createResp.Body)
	var createdTask entities.Task
	json.Unmarshal(body, &createdTask)

	// Now update the task
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
	err = json.Unmarshal(body, &updatedTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if updatedTask.Name != updateReq.Name {
		t.Errorf("Expected name '%s', got '%s'", updateReq.Name, updatedTask.Name)
	}
	if updatedTask.Status != updateReq.Status {
		t.Errorf("Expected status %d, got %d", updateReq.Status, updatedTask.Status)
	}
}

func TestSetupRoutes_UpdateTask_ValidationError(t *testing.T) {
	app := setupTestApp()

	// Try to update with invalid data (no need to create first for validation test)
	updateReq := requests.UpdateTaskRequest{Name: "", Status: 0}
	updateBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestSetupRoutes_DeleteTask(t *testing.T) {
	app := setupTestApp()

	// First create a task
	taskReq := requests.CreateTaskRequest{Name: "Delete Test Task", Status: 0}
	reqBody, _ := json.Marshal(taskReq)
	createReq := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := app.Test(createReq)

	body, _ := io.ReadAll(createResp.Body)
	var createdTask entities.Task
	json.Unmarshal(body, &createdTask)

	// Now delete the task
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/tasks/%d", createdTask.ID), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusNoContent {
		t.Errorf("Expected status %d, got %d", fiber.StatusNoContent, resp.StatusCode)
	}

	// Verify task is deleted
	getReq := httptest.NewRequest("GET", fmt.Sprintf("/tasks/%d", createdTask.ID), nil)
	getResp, _ := app.Test(getReq)
	if getResp.StatusCode != fiber.StatusBadRequest {
		t.Error("Expected task to be deleted")
	}
}

func TestSetupRoutes_DeleteTask_InvalidID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("DELETE", "/tasks/abc", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestSetupRoutes_DeleteTask_NotFound(t *testing.T) {
	app := setupTestApp()

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

func TestSetupRoutes_IntegrationFlow(t *testing.T) {
	app := setupTestApp()

	// 1. Check health
	healthReq := httptest.NewRequest("GET", "/health", nil)
	healthResp, _ := app.Test(healthReq)
	if healthResp.StatusCode != fiber.StatusOK {
		t.Error("Health check failed")
	}

	// 2. Get all tasks (should be empty)
	allReq := httptest.NewRequest("GET", "/tasks", nil)
	allResp, _ := app.Test(allReq)
	body, _ := io.ReadAll(allResp.Body)
	var tasks []entities.Task
	json.Unmarshal(body, &tasks)
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks initially, got %d", len(tasks))
	}

	// 3. Create a task
	taskReq := requests.CreateTaskRequest{Name: "Integration Task", Status: 0}
	reqBody, _ := json.Marshal(taskReq)
	createReq := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := app.Test(createReq)

	body, _ = io.ReadAll(createResp.Body)
	var createdTask entities.Task
	json.Unmarshal(body, &createdTask)

	// 4. Get all tasks (should have 1)
	allReq2 := httptest.NewRequest("GET", "/tasks", nil)
	allResp2, _ := app.Test(allReq2)
	body, _ = io.ReadAll(allResp2.Body)
	json.Unmarshal(body, &tasks)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task after creation, got %d", len(tasks))
	}

	// 5. Update the task
	updateReq := requests.UpdateTaskRequest{Name: "Updated Integration Task", Status: 1}
	updateBody, _ := json.Marshal(updateReq)
	putReq := httptest.NewRequest("PUT", fmt.Sprintf("/tasks/%d", createdTask.ID), bytes.NewBuffer(updateBody))
	putReq.Header.Set("Content-Type", "application/json")
	putResp, _ := app.Test(putReq)
	if putResp.StatusCode != fiber.StatusOK {
		t.Error("Task update failed")
	}

	// 6. Get the updated task
	getReq := httptest.NewRequest("GET", fmt.Sprintf("/tasks/%d", createdTask.ID), nil)
	getResp, _ := app.Test(getReq)
	body, _ = io.ReadAll(getResp.Body)
	var updatedTask entities.Task
	json.Unmarshal(body, &updatedTask)
	if updatedTask.Status != 1 {
		t.Errorf("Expected updated status 1, got %d", updatedTask.Status)
	}

	// 7. Delete the task
	deleteReq := httptest.NewRequest("DELETE", fmt.Sprintf("/tasks/%d", createdTask.ID), nil)
	deleteResp, _ := app.Test(deleteReq)
	if deleteResp.StatusCode != fiber.StatusNoContent {
		t.Error("Task deletion failed")
	}

	// 8. Verify task is gone
	getReq2 := httptest.NewRequest("GET", fmt.Sprintf("/tasks/%d", createdTask.ID), nil)
	getResp2, _ := app.Test(getReq2)
	if getResp2.StatusCode != fiber.StatusBadRequest {
		t.Error("Expected task to be deleted")
	}
}

func TestSetupRoutes_RouteRegistration(t *testing.T) {
	app := setupTestApp()

	// Test that all expected routes are registered by trying to access them
	routes := []struct {
		method string
		path   string
		body   string
	}{
		{"GET", "/health", ""},
		{"GET", "/tasks", ""},
		{"POST", "/tasks", `{"name":"test","status":0}`},
		{"GET", "/tasks/1", ""},
		{"PUT", "/tasks/1", `{"name":"test","status":0}`},
		{"DELETE", "/tasks/1", ""},
	}

	for _, route := range routes {
		var req *http.Request
		if route.body != "" {
			req = httptest.NewRequest(route.method, route.path, bytes.NewBuffer([]byte(route.body)))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req = httptest.NewRequest(route.method, route.path, nil)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Route %s %s failed: %v", route.method, route.path, err)
		}

		// We don't expect 404 (route not found) for any of these
		if resp.StatusCode == fiber.StatusNotFound {
			t.Errorf("Route %s %s not registered (got 404)", route.method, route.path)
		}
	}
}
