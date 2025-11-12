package request

import "time"

type CreateTaskRequest struct {
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Description string    `json:"description" validate:"max=1000"`
	Status      string    `json:"status" validate:"oneof=todo in_progress done"`
	Priority    string    `json:"priority" validate:"oneof=low medium high"`
	DueDate     time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"title" validate:"omitempty,min=1,max=255"`
	Description string    `json:"description" validate:"omitempty,max=1000"`
	Status      string    `json:"status" validate:"omitempty,oneof=todo in_progress done"`
	Priority    string    `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate     time.Time `json:"due_date"`
}
