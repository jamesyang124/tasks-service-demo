package entities

// Task represents a task entity with ID, name, and status.
type Task struct {
	ID     int    `json:"id"`                                     // Unique identifier for the task
	Name   string `json:"name" validate:"required,min=1,max=100"` // Task name (required, 1-100 chars)
	Status int    `json:"status" validate:"oneof=0 1"`            // Task status (0=incomplete, 1=complete)
}
