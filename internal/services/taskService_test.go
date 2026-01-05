package services

import (
	"context"
	"fmt"
	"math"
	"task-manager/internal/contextkeys"
	"task-manager/internal/models"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type mockTaskRepository struct {
	indexFn       func(uID int64) (models.TasksList, error)
	storeFn       func(ctx context.Context, p models.TaskPayload) error
	updateFn      func(ctx context.Context, p models.UpdateTask) (models.Task, error)
	showFn        func(id int) (models.Task, error)
	isTaskOwnerFn func(uID int64, id int) (bool, error)
	deleteTaskFn  func(id int) error
}

func (m mockTaskRepository) Store(ctx context.Context, p models.TaskPayload) error {
	if m.storeFn != nil {
		return m.storeFn(ctx, p)
	}
	return nil
}

func (m mockTaskRepository) Update(ctx context.Context, p models.UpdateTask) (models.Task, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, p)
	}
	return models.Task{}, nil
}

func (m mockTaskRepository) Show(id int) (models.Task, error) {
	if m.showFn != nil {
		return m.showFn(id)
	}
	return models.Task{}, nil
}

func (m mockTaskRepository) Index(uID int64) (models.TasksList, error) {
	if m.indexFn != nil {
		return m.indexFn(uID)
	}

	return models.TasksList{}, nil
}

func (m mockTaskRepository) Delete(id int) error {
	if m.deleteTaskFn != nil {
		return m.deleteTaskFn(id)
	}
	return nil
}

func (m mockTaskRepository) IsTaskOwner(uID int64, id int) (bool, error) {
	if m.isTaskOwnerFn != nil {
		return m.isTaskOwnerFn(uID, id)
	}
	return true, nil
}

func TestTaskService_GetTasksList(t *testing.T) {
	var tests = []struct {
		name            string
		mock            mockTaskRepository
		uID             int64
		expectedPayload models.TasksList
		expectsError    bool
		errorWanted     string
	}{
		{
			"valid ID, list with 1 item returned",
			mockTaskRepository{indexFn: func(uID int64) (models.TasksList, error) {
				return models.TasksList{
					Tasks: []models.Task{{ID: 1, Name: "Example task", Priority: models.PriorityLow, CreatedBy: 1}},
				}, nil
			}},
			1,
			models.TasksList{
				Tasks: []models.Task{{ID: 1, Name: "Example task", Priority: models.PriorityLow, CreatedBy: 1}},
			},
			false,
			"",
		},
		{
			"valid ID, list returned",
			mockTaskRepository{indexFn: func(uID int64) (models.TasksList, error) {
				return models.TasksList{
					Tasks: []models.Task{
						{ID: 1, Name: "Example task 1", Priority: models.PriorityLow, CreatedBy: 1},
						{ID: 2, Name: "Example task 2", Priority: models.PriorityMedium, CreatedBy: 1},
						{ID: 3, Name: "Example task 3", Priority: models.PriorityLow, CreatedBy: 1},
					},
				}, nil
			}},
			1,
			models.TasksList{
				Tasks: []models.Task{
					{ID: 1, Name: "Example task 1", Priority: models.PriorityLow, CreatedBy: 1},
					{ID: 2, Name: "Example task 2", Priority: models.PriorityMedium, CreatedBy: 1},
					{ID: 3, Name: "Example task 3", Priority: models.PriorityLow, CreatedBy: 1},
				},
			},
			false,
			"",
		},
		{
			"max ID value, no items returned",
			mockTaskRepository{indexFn: func(uID int64) (models.TasksList, error) {
				return models.TasksList{}, nil
			}},
			math.MaxInt64,
			models.TasksList{Tasks: nil},
			false,
			"",
		},
		{
			"valid ID, no items returned",
			mockTaskRepository{indexFn: func(uID int64) (models.TasksList, error) {
				return models.TasksList{}, nil
			}},
			2,
			models.TasksList{Tasks: nil},
			false,
			"",
		},
		{
			"invalid ID, no items returned",
			mockTaskRepository{},
			0,
			models.TasksList{
				Tasks: nil,
			},
			true,
			"GetTasksList: invalid user",
		},
		{
			"negative ID value, no items returned",
			mockTaskRepository{},
			-1,
			models.TasksList{
				Tasks: nil,
			},
			true,
			"GetTasksList: invalid user",
		},
		{
			"error while executing the query",
			mockTaskRepository{indexFn: func(uID int64) (models.TasksList, error) {
				return models.TasksList{}, fmt.Errorf("failed to execute query")
			}},
			12,
			models.TasksList{
				Tasks: nil,
			},
			true,
			"GetTasksList: failed to get data: failed to execute query",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewTaskService(tc.mock)
			tl, err := s.GetTasksList(tc.uID)
			if tc.expectsError && err == nil {
				t.Errorf("function is expected to return an error but it did not")
			}

			if !tc.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if !tc.expectsError {
				if diff := cmp.Diff(tc.expectedPayload, tl); diff != "" {
					t.Errorf("invalid data returned <-want, +got>\n%s", diff)
				}
			} else {
				if err.Error() != tc.errorWanted {
					t.Errorf("expected error <%s> but got <%s>", tc.errorWanted, err)
				}
			}
		})
	}
}

func TestTaskService_StoreTask(t *testing.T) {
	var tests = []struct {
		name         string
		mock         mockTaskRepository
		payload      models.TaskPayload
		expectsError bool
		errorWanted  string
	}{
		{
			"valid payload, low priority, task created",
			mockTaskRepository{storeFn: func(ctx context.Context, p models.TaskPayload) error {
				return nil
			}},
			models.TaskPayload{
				Name:        "LoremIpsum",
				Priority:    models.PriorityLow,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
			},
			false,
			"",
		},
		{
			"valid payload, medium priority, task created",
			mockTaskRepository{storeFn: func(ctx context.Context, p models.TaskPayload) error {
				return nil
			}},
			models.TaskPayload{
				Name:        "LoremIpsum",
				Priority:    models.PriorityMedium,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
			},
			false,
			"",
		},
		{
			"valid payload, high priority, task created",
			mockTaskRepository{storeFn: func(ctx context.Context, p models.TaskPayload) error {
				return nil
			}},
			models.TaskPayload{
				Name:        "LoremIpsum",
				Priority:    models.PriorityHigh,
				Description: "",
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, 12); return &t }(),
				CreatedAt:   nil,
			},
			false,
			"",
		},
		{
			"invalid payload, error returned",
			mockTaskRepository{storeFn: func(ctx context.Context, p models.TaskPayload) error {
				return fmt.Errorf("error executing the query")
			}},
			models.TaskPayload{
				Name:        "LoremIpsum",
				Priority:    21,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
			},
			true,
			"storeTask: error while storing the data: error executing the query",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewTaskService(tc.mock)
			err := s.StoreTask(context.Background(), tc.payload)
			if tc.expectsError && err == nil {
				t.Errorf("function was supposed to return an error but it did not")
			}

			if !tc.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if tc.expectsError && err != nil {
				if tc.errorWanted != err.Error() {
					t.Errorf("expected error <%s> but got <%s>", tc.errorWanted, err)
				}
			}
		})
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	var tests = []struct {
		name           string
		mock           mockTaskRepository
		payload        models.UpdateTask
		expectedResult models.Task
		expectsError   bool
		errorWanted    string
	}{
		{
			"name can be updated",
			mockTaskRepository{updateFn: func(ctx context.Context, p models.UpdateTask) (models.Task, error) {
				return models.Task{
					ID:          1,
					Name:        p.Name,
					Priority:    models.PriorityLow,
					Description: "",
					DueDate:     nil,
					CreatedAt:   nil,
					CreatedBy:   1}, nil
			}},
			models.UpdateTask{ID: 1, Name: "Lorem Ipsum"},
			models.Task{
				ID:          1,
				Name:        "Lorem Ipsum",
				Priority:    models.PriorityLow,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
				CreatedBy:   1,
			},
			false,
			"",
		},
		{
			"priority can be updated",
			mockTaskRepository{updateFn: func(ctx context.Context, p models.UpdateTask) (models.Task, error) {
				return models.Task{
					ID:          1,
					Name:        "LoremIpsum",
					Priority:    p.Priority,
					Description: "",
					DueDate:     nil,
					CreatedAt:   nil,
					CreatedBy:   1}, nil
			}},
			models.UpdateTask{ID: 1, Name: "LoremIpsum", Priority: models.PriorityMedium},
			models.Task{
				ID:          1,
				Name:        "LoremIpsum",
				Priority:    models.PriorityMedium,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
				CreatedBy:   1,
			},
			false,
			"",
		},
		{
			"description can be updated",
			mockTaskRepository{updateFn: func(ctx context.Context, p models.UpdateTask) (models.Task, error) {
				return models.Task{
					ID:          1,
					Name:        "LoremIpsum",
					Priority:    models.PriorityLow,
					Description: p.Description,
					DueDate:     nil,
					CreatedAt:   nil,
					CreatedBy:   1,
				}, nil
			}},
			models.UpdateTask{ID: 1, Description: "Lorem ipsum"},
			models.Task{
				ID:          1,
				Name:        "LoremIpsum",
				Priority:    models.PriorityLow,
				Description: "Lorem ipsum",
				DueDate:     nil,
				CreatedAt:   nil,
				CreatedBy:   1,
			},
			false,
			"",
		},
		{
			"description can be cleared",
			mockTaskRepository{updateFn: func(ctx context.Context, p models.UpdateTask) (models.Task, error) {
				return models.Task{
					ID:          1,
					Name:        "LoremIpsum",
					Priority:    models.PriorityLow,
					Description: p.Description,
					DueDate:     nil,
					CreatedAt:   nil,
					CreatedBy:   1,
				}, nil
			}},
			models.UpdateTask{
				ID:          1,
				Description: "",
			},
			models.Task{
				ID:          1,
				Name:        "LoremIpsum",
				Priority:    models.PriorityLow,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
				CreatedBy:   1,
			},
			false,
			"",
		},
		{
			"unauthorized, update failed",
			mockTaskRepository{updateFn: func(ctx context.Context, p models.UpdateTask) (models.Task, error) {
				uID := ctx.Value(contextkeys.UserID).(int64)
				task := models.Task{
					ID:          1,
					Name:        "LoremIpsum",
					Priority:    models.PriorityLow,
					Description: "",
					DueDate:     nil,
					CreatedAt:   nil,
					CreatedBy:   1,
				}

				if uID != task.CreatedBy {
					return task, fmt.Errorf("user not authorized for this action")
				}

				return task, nil
			}},
			models.UpdateTask{ID: 1, Description: "Lorem ipsum"},
			models.Task{
				ID:          1,
				Name:        "LoremIpsum",
				Priority:    models.PriorityLow,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
				CreatedBy:   1,
			},
			true,
			"UpdateTask: user not authorized for this action",
		},
		{
			"invalid task ID",
			mockTaskRepository{updateFn: func(ctx context.Context, p models.UpdateTask) (models.Task, error) {
				return models.Task{}, fmt.Errorf("update: failed to get task from db")
			}},
			models.UpdateTask{ID: 2, Description: "Lorem ipsum"},
			models.Task{
				ID:          1,
				Name:        "LoremIpsum",
				Priority:    models.PriorityLow,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
				CreatedBy:   1,
			},
			true,
			"UpdateTask: update: failed to get task from db",
		},
		{
			"failed to update the task",
			mockTaskRepository{updateFn: func(ctx context.Context, p models.UpdateTask) (models.Task, error) {
				return models.Task{}, fmt.Errorf("update: failed to execute the update query")
			}},
			models.UpdateTask{ID: 1, Description: "Lorem ipsum"},
			models.Task{
				ID:          1,
				Name:        "LoremIpsum",
				Priority:    models.PriorityLow,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
				CreatedBy:   1,
			},
			true,
			"UpdateTask: update: failed to execute the update query",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewTaskService(tc.mock)
			ctx := context.WithValue(context.Background(), contextkeys.UserID, int64(2))

			task, err := s.UpdateTask(ctx, tc.payload)

			if tc.expectsError && err == nil {
				t.Errorf("function was supposed to return an error but it did not")
			}

			if !tc.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if !tc.expectsError {
				if diff := cmp.Diff(tc.expectedResult, task); diff != "" {
					t.Errorf("invalid data returned, <-want, +got>\n%s", diff)
				}
			}
			if tc.expectsError && err != nil {
				if err.Error() != tc.errorWanted {
					t.Errorf("expected error <%s> but got <%s>", tc.errorWanted, err)
				}
			}
		})
	}
}

func TestTaskService_ShowTask(t *testing.T) {
	var tests = []struct {
		name           string
		mock           mockTaskRepository
		taskID         int
		expectedResult models.Task
		expectsError   bool
		errorWanted    string
	}{
		{
			"correct id, task returned",
			mockTaskRepository{showFn: func(id int) (models.Task, error) {
				task := models.Task{
					ID:          1,
					Name:        "Lorem Ipsum",
					Priority:    models.PriorityLow,
					Description: "",
					DueDate:     nil,
					CreatedAt:   nil,
					CreatedBy:   1,
				}
				if task.ID == id {
					return task, nil
				}

				return models.Task{}, fmt.Errorf("error")
			}},
			1,
			models.Task{
				ID:          1,
				Name:        "Lorem Ipsum",
				Priority:    models.PriorityLow,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
				CreatedBy:   1,
			},
			false,
			"",
		},
		{
			"valid ID, no results",
			mockTaskRepository{showFn: func(id int) (models.Task, error) {
				return models.Task{}, fmt.Errorf("no results for given ID")
			}},
			123,
			models.Task{
				ID:          1,
				Name:        "Lorem Ipsum",
				Priority:    models.PriorityLow,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
				CreatedBy:   1,
			},
			true,
			"ShowTask: no results for given ID",
		},
		{
			"invalid ID, error returned",
			mockTaskRepository{showFn: func(id int) (models.Task, error) {
				return models.Task{}, fmt.Errorf("no results for given ID")
			}},
			-1,
			models.Task{},
			true,
			"ShowTask: no results for given ID",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewTaskService(tc.mock)

			task, err := s.ShowTask(tc.taskID)
			if tc.expectsError && err == nil {
				t.Errorf("function is supposed to return an error but it did not")
			}

			if !tc.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if !tc.expectsError {
				if diff := cmp.Diff(tc.expectedResult, task); diff != "" {
					t.Errorf("invalid data returned <-want, +got>\n%s", diff)
				}
			}

			if tc.expectsError && err != nil {
				if tc.errorWanted != err.Error() {
					t.Errorf("expected error <%s> but got <%s>", tc.errorWanted, err)
				}
			}
		})
	}
}

func TestTaskService_IsTaskOwner(t *testing.T) {
	var tests = []struct {
		name               string
		mock               mockTaskRepository
		uID                int64
		tID                int
		expectedIsOwnerVal bool
		expectsError       bool
		errorWanted        string
	}{
		{
			"valid values, user is owner",
			mockTaskRepository{isTaskOwnerFn: func(uID int64, id int) (bool, error) {
				return true, nil
			}},
			1,
			1,
			true,
			false,
			"",
		},
		{
			"valid values, user is not owner",
			mockTaskRepository{isTaskOwnerFn: func(uID int64, id int) (bool, error) {
				return false, nil
			}},
			1,
			2,
			false,
			false,
			"",
		},
		{
			"invalid user ID value",
			mockTaskRepository{isTaskOwnerFn: func(uID int64, id int) (bool, error) {
				return false, fmt.Errorf("invalid parameter passed")
			}},
			-1,
			2,
			false,
			true,
			"failed to check if user is a task owner: invalid parameter passed",
		},
		{
			"invalid task ID value",
			mockTaskRepository{isTaskOwnerFn: func(uID int64, id int) (bool, error) {
				return false, fmt.Errorf("invalid parameter passed")
			}},
			1,
			-2,
			false,
			true,
			"failed to check if user is a task owner: invalid parameter passed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewTaskService(tc.mock)

			isOwner, err := s.IsTaskOwner(tc.uID, tc.tID)

			if tc.expectsError && err == nil {
				t.Errorf("function is supposed to return an error but it did not")
			}

			if !tc.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if tc.expectsError && err != nil {
				if tc.errorWanted != err.Error() {
					t.Errorf("expected error <%s> but got <%s>", tc.errorWanted, err)
				}
			}

			if tc.expectedIsOwnerVal != isOwner {
				t.Errorf("function is supposed to return %t but got %t instead", tc.expectedIsOwnerVal, isOwner)
			}
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	var tests = []struct {
		name         string
		mock         mockTaskRepository
		tID          int
		uID          int64
		expectsError bool
		errorWanted  string
	}{
		{
			"valid values, no error",
			mockTaskRepository{deleteTaskFn: func(id int) error {
				return nil
			}},
			1,
			1,
			false,
			"",
		},
		{
			"valid values, user not an owner",
			mockTaskRepository{deleteTaskFn: func(id int) error {
				return fmt.Errorf("you are not authorized to execute this action")
			}},
			1,
			12,
			true,
			"DeleteTask: you are not authorized to execute this action",
		},
		{
			"valid values, failed to delete",
			mockTaskRepository{deleteTaskFn: func(id int) error {
				return fmt.Errorf("delete: failed to execute query")
			}},
			1,
			12,
			true,
			"DeleteTask: delete: failed to execute query",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewTaskService(tc.mock)

			err := s.DeleteTask(tc.tID, tc.uID)

			if !tc.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if tc.expectsError && err == nil {
				t.Errorf("function is supposed to return an error but it did not")
			}

			if tc.expectsError && err != nil {
				if tc.errorWanted != err.Error() {
					t.Errorf("expected error <%s> but got <%s>", tc.errorWanted, err)
				}
			}
		})
	}
}
