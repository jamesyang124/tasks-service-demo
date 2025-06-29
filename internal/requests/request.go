package requests

import apperrors "tasks-service-demo/internal/errors"

// Package requests defines request types and validation logic for the Task API.

// CreateTaskRequest represents the request body for creating a task.
type CreateTaskRequest struct {
	Name   string `json:"name" validate:"required,min=1,max=100"`
	Status int    `json:"status" validate:"oneof=0 1"`
}

// UpdateTaskRequest represents the request body for updating a task.
type UpdateTaskRequest struct {
	Name   string `json:"name" validate:"required,min=1,max=100"`
	Status int    `json:"status" validate:"oneof=0 1"`
}

// Validatable is an interface for request validation.
type Validatable interface {
	Validate() *apperrors.AppError
}

// Validate validates the CreateTaskRequest fields.
func (c CreateTaskRequest) Validate() *apperrors.AppError {
	return ValidateStruct(&c)
}

// Validate validates the UpdateTaskRequest fields.
func (u UpdateTaskRequest) Validate() *apperrors.AppError {
	return ValidateStruct(&u)
}
