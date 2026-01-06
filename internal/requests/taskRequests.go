package requests

import (
	"strings"
	"task-manager/internal/models"
	"time"
)

type CreateTasksRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Priority    models.Priority `json:"priority"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
	CreatedAt   *time.Time      `json:"created-at,omitempty"`
}

func (r CreateTasksRequest) Validate() ValidationResult {
	var res ValidationResult
	res.Validated = true
	l := len(strings.TrimSpace(r.Name))
	if l < 3 || l > 64 {
		res.SetFailed("name must be between 3 and 64 characters")
	}

	if r.Priority < models.PriorityLow || r.Priority > models.PriorityHigh {
		res.SetFailed("invalid priority value")
	}

	if r.Priority == models.PriorityHigh && r.DueDate == nil {
		res.SetFailed("for priority high due date is required")
	}

	if len(r.Description) > 255 {
		res.SetFailed("description too long")
	}

	return res
}

type UpdateTaskRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Priority    models.Priority `json:"priority"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
	CreatedAt   *time.Time      `json:"created-at,omitempty"`
}

func (r UpdateTaskRequest) Validate() ValidationResult {
	var res ValidationResult
	res.Validated = true
	l := len(strings.TrimSpace(r.Name))
	if l < 3 || l > 64 {
		res.SetFailed("name must be between 3 and 64 characters")
	}

	if r.Priority < models.PriorityLow || r.Priority > models.PriorityHigh {
		res.SetFailed("invalid priority value")
	}

	if r.Priority == models.PriorityHigh && r.DueDate == nil {
		res.SetFailed("for priority high due date is required")
	}

	if len(r.Description) > 255 {
		res.SetFailed("description too long")
	}

	return res
}
