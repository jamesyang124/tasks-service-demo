package models

type Validatable interface {
	Validate() error
}

type CreateTaskRequest struct {
	Name   string `json:"name" validate:"required,min=1,max=100"`
	Status int    `json:"status" validate:"oneof=0 1"`
}

func (c CreateTaskRequest) Validate() error {
	return ValidateStruct(&c)
}

type UpdateTaskRequest struct {
	Name   string `json:"name" validate:"required,min=1,max=100"`
	Status int    `json:"status" validate:"oneof=0 1"`
}

func (u UpdateTaskRequest) Validate() error {
	return ValidateStruct(&u)
}
