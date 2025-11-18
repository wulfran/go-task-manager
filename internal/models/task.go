package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

func (p *Priority) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch strings.ToLower(s) {
		case "low":
			*p = PriorityLow
			return nil
		case "medium":
			*p = PriorityMedium
			return nil
		case "high":
			*p = PriorityHigh
			return nil
		default:
			return fmt.Errorf("invalid priority value: %s", data)
		}
	}

	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		switch i {
		case int(PriorityLow), int(PriorityMedium), int(PriorityHigh):
			*p = Priority(i)
			return nil
		default:
			return fmt.Errorf("invalid priority value: %d", i)
		}
	}

	return fmt.Errorf("invalid payload: %v", data)
}

type Task struct {
	ID          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Priority    Priority   `json:"priority" db:"priority"`
	Description string     `json:"description,omitempty" db:"description"`
	DueDate     *time.Time `json:"due_date" db:"due_date"`
	CreatedAt   *time.Time `json:"created_at" db:"created_at"`
	CreatedBy   int64      `json:"created_by" db:"created_by"`
}

type TasksList struct {
	Tasks []Task `json:"tasks"`
}

type TaskPayload struct {
	Name        string     `json:"name"`
	Priority    Priority   `json:"priority"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

type UpdateTask struct {
	ID          int64
	Name        string     `json:"name,omitempty"`
	Priority    Priority   `json:"priority,omitempty"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}
