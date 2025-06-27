package services

import (
	"tasks-service-demo/internal/models"
	"tasks-service-demo/internal/storage"
)

type TaskService struct{}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (s *TaskService) store() storage.Store {
	return storage.GetStore()
}

func (s *TaskService) GetAllTasks() []*models.Task {
	return s.store().GetAll()
}

func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
	return s.store().GetByID(id)
}

func (s *TaskService) CreateTask(req *models.CreateTaskRequest) (*models.Task, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	task := &models.Task{
		Name:   req.Name,
		Status: req.Status,
	}
	
	if err := s.store().Create(task); err != nil {
		return nil, err
	}
	
	return task, nil
}

func (s *TaskService) UpdateTask(id int, req *models.UpdateTaskRequest) (*models.Task, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	task := &models.Task{
		Name:   req.Name,
		Status: req.Status,
	}
	
	if err := s.store().Update(id, task); err != nil {
		return nil, err
	}
	
	return task, nil
}

func (s *TaskService) DeleteTask(id int) error {
	return s.store().Delete(id)
}