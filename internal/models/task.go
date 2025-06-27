package models

import "encoding/json"

type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name" validate:"required,min=1,max=100"`
	Status int    `json:"status" validate:"oneof=0 1"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func (e *ErrorResponse) ToJSON() []byte {
	data, _ := json.Marshal(e)
	return data
}
