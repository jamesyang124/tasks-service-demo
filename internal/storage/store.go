package storage

import (
	"sync"
	"tasks-service-demo/internal/entities"
	apperrors "tasks-service-demo/internal/errors"
)

// Store defines the interface for all storage implementations
type Store interface {
	Create(task *entities.Task) *apperrors.AppError         // Creates a new task
	GetByID(id int) (*entities.Task, *apperrors.AppError)   // Retrieves a task by ID
	GetAll() []*entities.Task                               // Retrieves all tasks
	Update(id int, task *entities.Task) *apperrors.AppError // Updates an existing task
	Delete(id int) *apperrors.AppError                      // Deletes a task by ID
}

// Singleton pattern for application-wide store instance
var (
	instance Store
	once     sync.Once
)

// InitStore initializes the singleton store instance
func InitStore(store Store) {
	once.Do(func() {
		instance = store
	})
}

// GetStore returns the singleton store instance
func GetStore() Store {
	return instance
}

// ResetStore reset store instance
func ResetStore() {
	instance = nil
	once = sync.Once{}
}
