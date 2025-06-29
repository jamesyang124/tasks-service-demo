package naive

import (
	"sync"
	"tasks-service-demo/internal/entities"

	apperrors "tasks-service-demo/internal/errors"
)

// MemoryStore provides an in-memory storage implementation using a map and mutex
type MemoryStore struct {
	tasks  map[int]*entities.Task // Map to store tasks by ID
	mu     sync.RWMutex           // Read-write mutex for thread safety
	nextID int                    // Auto-incrementing ID counter
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:  make(map[int]*entities.Task),
		nextID: 1,
	}
}

// Create stores a new task with an auto-generated ID
func (s *MemoryStore) Create(task *entities.Task) *apperrors.AppError {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.nextID
	s.nextID++
	s.tasks[task.ID] = task
	return nil
}

// GetByID retrieves a task by its ID, returns error if not found
func (s *MemoryStore) GetByID(id int) (*entities.Task, *apperrors.AppError) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, apperrors.ErrTaskNotFound
	}
	return task, nil
}

// GetAll returns all tasks in the store
func (s *MemoryStore) GetAll() []*entities.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*entities.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// Update modifies an existing task by ID, returns error if not found
func (s *MemoryStore) Update(id int, updatedTask *entities.Task) *apperrors.AppError {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return apperrors.ErrTaskNotFound
	}

	updatedTask.ID = id
	s.tasks[id] = updatedTask
	return nil
}

// Delete removes a task by ID, returns error if not found
func (s *MemoryStore) Delete(id int) *apperrors.AppError {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return apperrors.ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}
