package storage

import (
	"sync"
	"tasks-service-demo/internal/models"
)

type Store interface {
	Create(task *models.Task) error
	GetByID(id int) (*models.Task, error)
	GetAll() []*models.Task
	Update(id int, task *models.Task) error
	Delete(id int) error
}

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
	if instance == nil {
		// Default to memory store if not initialized
		InitStore(NewMemoryStore())
	}
	return instance
}

// ResetStore resets the singleton for testing purposes
func ResetStore() {
	instance = nil
	once = sync.Once{}
}
