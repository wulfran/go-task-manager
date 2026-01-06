package requests

import (
	"math"
	"strings"
	"task-manager/internal/models"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type taskRequestTestCase struct {
	name        string
	taskName    string
	description string
	priority    models.Priority
	dueDate     *time.Time
	expected    ValidationResult
}

func testCases() []taskRequestTestCase {
	return []taskRequestTestCase{
		{
			name:        "correct struct, no errors",
			taskName:    "Lorem ipsum",
			description: "Dolor Et",
			priority:    models.PriorityLow,
			expected: ValidationResult{
				Validated: true,
				Message:   "",
			},
		},
		{
			name:        "priority high with due date, no errors",
			taskName:    "Lorem ipsum",
			description: "Dolor Et",
			priority:    models.PriorityHigh,
			dueDate: func() *time.Time {
				t := time.Date(2027, 12, 12, 0, 0, 0, 0, time.UTC)
				return &t
			}(),
			expected: ValidationResult{
				Validated: true,
				Message:   "",
			},
		},
		{
			name:        "missing name",
			taskName:    "",
			description: "Dolor Et",
			priority:    models.PriorityLow,
			expected: ValidationResult{
				Validated: false,
				Message:   "name must be between 3 and 64 characters",
			},
		},
		{
			name:        "name out of allowed length",
			taskName:    strings.Repeat("a", 65),
			description: "Dolor Et",
			priority:    models.PriorityLow,
			expected: ValidationResult{
				Validated: false,
				Message:   "name must be between 3 and 64 characters",
			},
		},
		{
			name:        "name all whitespaces",
			taskName:    "    ",
			description: "Dolor Et",
			priority:    models.PriorityLow,
			expected: ValidationResult{
				Validated: false,
				Message:   "name must be between 3 and 64 characters",
			},
		},
		{
			name:        "invalid priority value",
			taskName:    "Lorem Ipsum",
			description: "Dolor Et",
			priority:    math.MaxInt,
			expected: ValidationResult{
				Validated: false,
				Message:   "invalid priority value",
			},
		},
		{
			name:        "priority negative value",
			taskName:    "Lorem Ipsum",
			description: "Dolor Et",
			priority:    -1,
			expected: ValidationResult{
				Validated: false,
				Message:   "invalid priority value",
			},
		},
		{
			name:        "priority high, missing due date",
			taskName:    "Lorem Ipsum",
			description: "Dolor Et",
			priority:    models.PriorityHigh,
			expected: ValidationResult{
				Validated: false,
				Message:   "for priority high due date is required",
			},
		},
		{
			name:        "description too long",
			taskName:    "Lorem Ipsum",
			description: strings.Repeat("a", 256),
			priority:    models.PriorityLow,
			expected: ValidationResult{
				Validated: false,
				Message:   "description too long",
			},
		},
		{
			name:     "multiple error messages",
			taskName: "",
			priority: math.MaxInt,
			expected: ValidationResult{
				Validated: false,
				Message:   "name must be between 3 and 64 characters, invalid priority value",
			},
		},
	}
}

func TestCreateTasksRequest_Validate(t *testing.T) {
	tests := testCases()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctr := CreateTasksRequest{
				Name:        tc.taskName,
				Description: tc.description,
				Priority:    tc.priority,
				DueDate:     tc.dueDate,
			}

			result := ctr.Validate()

			if diff := cmp.Diff(tc.expected, result); diff != "" {
				t.Errorf("unexpected validation result, <-want, +got>\n%s", diff)
			}
		})
	}
}

func TestUpdateTaskRequest_Validate(t *testing.T) {
	tests := testCases()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			utr := UpdateTaskRequest{
				Name:        tc.taskName,
				Description: tc.description,
				Priority:    tc.priority,
				DueDate:     tc.dueDate,
			}
			result := utr.Validate()

			if diff := cmp.Diff(tc.expected, result); diff != "" {
				t.Errorf("unexpected validation result, <-want, +got>\n%s", diff)
			}
		})
	}
}
