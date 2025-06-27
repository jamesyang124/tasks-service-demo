package storage

import (
	"errors"
	"sync"
	"tasks-service-demo/internal/models"
)

type MemoryStore struct {
	tasks  map[int]*models.Task
	mu     sync.RWMutex
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:  make(map[int]*models.Task),
		nextID: 1,
	}
}

func (s *MemoryStore) Create(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.nextID
	s.nextID++
	s.tasks[task.ID] = task
	return nil
}

func (s *MemoryStore) GetByID(id int) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (s *MemoryStore) GetAll() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *MemoryStore) Update(id int, updatedTask *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return errors.New("task not found")
	}

	updatedTask.ID = id
	s.tasks[id] = updatedTask
	return nil
}

func (s *MemoryStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(s.tasks, id)
	return nil
}
