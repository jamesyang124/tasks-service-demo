package services

import (
	"tasks-service-demo/internal/entities"
	apperrors "tasks-service-demo/internal/errors"
	"tasks-service-demo/internal/logger"
	"tasks-service-demo/internal/requests"
	"tasks-service-demo/internal/storage"
)

// Package services implements business logic for the Task API.

// TaskService provides methods for managing tasks.
type TaskService struct{}

// NewTaskService creates a new TaskService instance.
func NewTaskService() *TaskService {
	return &TaskService{}
}

func (s *TaskService) store() storage.Store {
	return storage.GetStore()
}

// GetAllTasks returns all tasks from the store.
func (s *TaskService) GetAllTasks() []*entities.Task {
	return s.store().GetAll()
}

// GetTaskByID returns a task by its ID, or an error if not found.
func (s *TaskService) GetTaskByID(id int) (*entities.Task, *apperrors.AppError) {
	task, err := s.store().GetByID(id)
	if err != nil {
		logger.Get().Error(err)
		return nil, err
	}
	return task, nil
}

// CreateTask creates a new task from the given request.
func (s *TaskService) CreateTask(req *requests.CreateTaskRequest) (*entities.Task, *apperrors.AppError) {
	task := &entities.Task{
		Name:   req.Name,
		Status: req.Status,
	}

	if err := s.store().Create(task); err != nil {
		logger.Get().Error(err)
		return nil, err
	}

	return task, nil
}

// UpdateTask updates an existing task by ID with the given request.
func (s *TaskService) UpdateTask(id int, req *requests.UpdateTaskRequest) (*entities.Task, *apperrors.AppError) {
	task := &entities.Task{
		Name:   req.Name,
		Status: req.Status,
	}

	if err := s.store().Update(id, task); err != nil {
		logger.Get().Error(err)
		return nil, err
	}

	return task, nil
}

// DeleteTask deletes a task by its ID. Returns nil if not found (idempotent).
func (s *TaskService) DeleteTask(id int) *apperrors.AppError {
	err := s.store().Delete(id)
	if err != nil {
		// RESTful design: DELETE should be idempotent
		logger.Get().Error(err)
	}
	return nil
}
