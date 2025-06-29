package storage

import (
	"sync"
	"tasks-service-demo/internal/models"
)

// Store defines the interface for all storage implementations
type Store interface {
	Create(task *models.Task) error
	GetByID(id int) (*models.Task, error)
	GetAll() []*models.Task
	Update(id int, task *models.Task) error
	Delete(id int) error
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
}
