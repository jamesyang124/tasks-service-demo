package xsync

import (
	"sync/atomic"
	"tasks-service-demo/internal/entities"

	apperrors "tasks-service-demo/internal/errors"
	"github.com/puzpuzpuz/xsync/v3"
)

// XSyncStore provides an in-memory storage implementation using xsync.Map
type XSyncStore struct {
	tasks  *xsync.MapOf[int, *entities.Task] // Concurrent map to store tasks by ID
	nextID int64                             // Atomic counter for ID generation
}

func NewXSyncStore() *XSyncStore {
	return &XSyncStore{
		tasks:  xsync.NewMapOf[int, *entities.Task](),
		nextID: 1,
	}
}

// Create stores a new task with an auto-generated ID
func (s *XSyncStore) Create(task *entities.Task) *apperrors.AppError {
	// Generate unique ID atomically
	id := int(atomic.AddInt64(&s.nextID, 1) - 1)
	task.ID = id
	
	s.tasks.Store(id, task)
	return nil
}

// GetByID retrieves a task by its ID, returns error if not found
func (s *XSyncStore) GetByID(id int) (*entities.Task, *apperrors.AppError) {
	task, ok := s.tasks.Load(id)
	if !ok {
		return nil, apperrors.ErrTaskNotFound
	}
	return task, nil
}

// GetAll returns all tasks in the store
func (s *XSyncStore) GetAll() []*entities.Task {
	tasks := make([]*entities.Task, 0)
	
	s.tasks.Range(func(key int, value *entities.Task) bool {
		tasks = append(tasks, value)
		return true // Continue iteration
	})
	
	return tasks
}

// Update modifies an existing task by ID, returns error if not found
func (s *XSyncStore) Update(id int, updatedTask *entities.Task) *apperrors.AppError {
	// Check if task exists first
	if _, ok := s.tasks.Load(id); !ok {
		return apperrors.ErrTaskNotFound
	}
	
	updatedTask.ID = id
	s.tasks.Store(id, updatedTask)
	return nil
}

// Delete removes a task by ID, returns error if not found
func (s *XSyncStore) Delete(id int) *apperrors.AppError {
	// Check if task exists first
	if _, ok := s.tasks.Load(id); !ok {
		return apperrors.ErrTaskNotFound
	}
	
	s.tasks.Delete(id)
	return nil
}