package repository

import (
	"context"
	"log"
	"strconv"
	"task-manager/internal/contextkeys"
	"task-manager/internal/models"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func storeTask(t *testing.T, r TaskRepository, p models.TaskPayload, ctx context.Context) error {
	t.Helper()
	err := r.Store(ctx, p)

	return err
}

func TestTaskRepository_Store(t *testing.T) {
	taskRepo := NewTaskRepository(*testDB)
	var tests = []struct {
		name             string
		payload          models.TaskPayload
		shouldCreateUser bool
		expectsError     bool
		errorWanted      string
	}{
		{
			"valid payload",
			models.TaskPayload{
				Name:     "Lorem",
				Priority: models.PriorityLow,
			},
			true,
			false,
			"",
		},
		{
			"no user in context",
			models.TaskPayload{
				Name:     "Lorem",
				Priority: models.PriorityLow,
			},
			false,
			true,
			"store: failed to get user id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, _ = testDB.Exec("TRUNCATE tasks, users RESTART IDENTITY CASCADE")
			})

			ctx := context.Background()
			if tc.shouldCreateUser {
				testCreateUser(t, *testDB)
				ctx = context.WithValue(ctx, contextkeys.UserID, int64(1))
			}
			err := storeTask(t, taskRepo, tc.payload, ctx)

			if tc.expectsError {
				if err == nil {
					t.Error("function was supposed to return an error but it did not")
				}
				if tc.errorWanted != err.Error() {
					t.Errorf("wrong error returned, expected <%s> but got <%s>", tc.errorWanted, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
			}
		})
	}
}

func TestTaskRepository_Update(t *testing.T) {
	taskRepo := NewTaskRepository(*testDB)
	var tests = []struct {
		name           string
		initialTask    models.TaskPayload
		payload        models.UpdateTask
		expectedResult models.Task
		userID         int64
		ownerID        int64
		expectsError   bool
		errorWanted    string
	}{
		{
			"valid payload",
			models.TaskPayload{
				Name:      "Lorem Ipsum",
				Priority:  models.PriorityLow,
				DueDate:   nil,
				CreatedAt: nil,
			},
			models.UpdateTask{
				ID:       1,
				Name:     "Lorem Ipsum Dolor Et",
				Priority: models.PriorityMedium,
			},
			models.Task{
				ID:          1,
				Name:        "Lorem Ipsum Dolor Et",
				Priority:    models.PriorityMedium,
				Description: "",
				DueDate:     nil,
				CreatedAt:   nil,
				CreatedBy:   1,
			},
			1,
			1,
			false,
			"",
		},
		{
			"non existing task",
			models.TaskPayload{
				Name:      "Lorem Ipsum",
				Priority:  models.PriorityLow,
				DueDate:   nil,
				CreatedAt: nil,
			},
			models.UpdateTask{
				ID:       2,
				Name:     "Lorem Ipsum Dolor Et",
				Priority: models.PriorityMedium,
			},
			models.Task{},
			1,
			1,
			true,
			"update: failed to get task from db: sql: no rows in result set",
		},
		{
			"wrong userID in context",
			models.TaskPayload{
				Name:      "Lorem Ipsum",
				Priority:  models.PriorityLow,
				DueDate:   nil,
				CreatedAt: nil,
			},
			models.UpdateTask{
				ID:       1,
				Name:     "Lorem Ipsum Dolor Et",
				Priority: models.PriorityMedium,
			},
			models.Task{},
			2,
			1,
			true,
			"user not authorized for this action",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, _ = testDB.Exec("TRUNCATE tasks, users RESTART IDENTITY CASCADE")
			})
			ctx := context.Background()
			testCreateUser(t, *testDB)
			ctx = context.WithValue(ctx, contextkeys.UserID, tc.ownerID)

			_ = storeTask(t, taskRepo, tc.initialTask, ctx)

			ctx = context.WithValue(ctx, contextkeys.UserID, tc.userID)

			task, err := taskRepo.Update(ctx, tc.payload)
			if tc.expectsError {
				if err == nil {
					t.Error("function was supposed to return an error but it did not")
				}
				if tc.errorWanted != err.Error() {
					t.Errorf("wrong error returned, expected <%s> but got <%s>", tc.errorWanted, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}

				if diff := cmp.Diff(tc.expectedResult, task); diff != "" {
					t.Errorf("wrong data returned: <-want, +got>\n%s", diff)
				}
			}
		})
	}
}

func TestTaskRepository_Show(t *testing.T) {
	var tests = []struct {
		name             string
		shouldCreateTask bool
		taskID           int
		expectedResult   models.Task
		expectsError     bool
		errorWanted      string
	}{
		{
			"valid id",
			true,
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
			"invalid id",
			false,
			123,
			models.Task{},
			true,
			"no results for given ID",
		},
		{
			"negative id",
			false,
			-1,
			models.Task{},
			true,
			"no results for given ID",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, _ = testDB.Exec("TRUNCATE tasks, users RESTART IDENTITY CASCADE ")
			})

			ctx := context.Background()
			taskRepo := NewTaskRepository(*testDB)
			if tc.shouldCreateTask {
				testCreateUser(t, *testDB)
				p := models.TaskPayload{
					Name:     "Lorem Ipsum",
					Priority: models.PriorityLow,
				}
				ctx = context.WithValue(ctx, contextkeys.UserID, int64(1))

				_ = storeTask(t, taskRepo, p, ctx)
			}

			task, err := taskRepo.Show(tc.taskID)
			if tc.expectsError {
				if err == nil {
					t.Error("function was supposed to return an error but it did not")
				}
				if tc.errorWanted != err.Error() {
					t.Errorf("wrong error returned, expected <%s> but got <%s>", tc.errorWanted, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}

				if diff := cmp.Diff(tc.expectedResult, task); diff != "" {
					t.Errorf("wrong data returned, <-want, +got>\n%s", diff)
				}
			}
		})
	}
}

func TestTaskRepository_Index(t *testing.T) {
	var tests = []struct {
		name              string
		userID            int64
		shouldCreateTask  bool
		howManyTasks      int
		expectedTaskCount int
		expectedResult    models.TasksList
		expectsError      bool
		errorWanted       string
	}{
		{
			"valid id one item",
			1,
			true,
			1,
			1,
			models.TasksList{
				Tasks: []models.Task{
					{
						1,
						"Lorem Ipsum1",
						models.PriorityLow,
						"",
						nil,
						nil,
						1,
					},
				},
			},
			false,
			"",
		},
		{
			"valid id more items",
			1,
			true,
			2,
			2,
			models.TasksList{
				Tasks: []models.Task{
					{
						1,
						"Lorem Ipsum1",
						models.PriorityLow,
						"",
						nil,
						nil,
						1,
					},
					{
						2,
						"Lorem Ipsum12",
						models.PriorityLow,
						"",
						nil,
						nil,
						1,
					},
				},
			},
			false,
			"",
		},
		{
			"non-existing userID",
			2,
			true,
			1,
			0,
			models.TasksList{},
			false,
			"",
		},
		{
			"valid id no items",
			1,
			false,
			0,
			0,
			models.TasksList{
				Tasks: nil,
			},
			false,
			"",
		},
		{
			"invalid id no items",
			-1,
			false,
			0,
			0,
			models.TasksList{
				Tasks: nil,
			},
			false,
			"",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, _ = testDB.Exec("TRUNCATE tasks, users RESTART IDENTITY CASCADE ")
			})

			taskRepo := NewTaskRepository(*testDB)
			ctx := context.Background()

			if tc.shouldCreateTask {
				p := models.TaskPayload{
					Name:     "Lorem Ipsum",
					Priority: models.PriorityLow,
				}
				testCreateUser(t, *testDB)
				ctx = context.WithValue(ctx, contextkeys.UserID, int64(1))

				for i := 1; i <= tc.howManyTasks; i++ {
					p.Name += strconv.Itoa(i)
					_ = storeTask(t, taskRepo, p, ctx)
				}
			}

			taskList, err := taskRepo.Index(tc.userID)
			if tc.expectsError {
				if err == nil {
					t.Error("function was supposed to return an error but it did not")
				}
				if tc.errorWanted != err.Error() {
					t.Errorf("unexpected error returned, wanted <%s> got <%s>", tc.errorWanted, err)
				}
			} else {
				if len(taskList.Tasks) != tc.expectedTaskCount {
					t.Errorf("wrong tasks count, wanted %d got %d", tc.expectedTaskCount, len(taskList.Tasks))
				}

				if diff := cmp.Diff(tc.expectedResult, taskList); diff != "" {
					t.Errorf("wrong data returned, <-want, +got>\n%s", diff)
				}
			}
		})
	}
}

func TestTaskRepository_Delete(t *testing.T) {
	var tests = []struct {
		name            string
		taskID          int
		shouldStoreTask bool
	}{
		{
			"Valid ID, delete task",
			1,
			true,
		},
		{
			"no tasks in db",
			0,
			false,
		},
		{
			"invalid task passed",
			-1,
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, _ = testDB.Exec("TRUNCATE tasks, users RESTART IDENTITY CASCADE ")
			})

			ctx := context.Background()
			taskRepo := NewTaskRepository(*testDB)

			if tc.shouldStoreTask {
				p := models.TaskPayload{
					Name:     "Lorem Ipsum",
					Priority: models.PriorityLow,
				}

				testCreateUser(t, *testDB)
				ctx := context.WithValue(ctx, contextkeys.UserID, int64(1))

				_ = storeTask(t, taskRepo, p, ctx)
			}

			err := taskRepo.Delete(tc.taskID)

			log.Println(err)
		})
	}
}

func TestTaskRepository_IsTaskOwner(t *testing.T) {
	var tests = []struct {
		name           string
		taskID         int
		userID         int64
		expectedResult bool
	}{
		{
			"valid data, user is owner",
			1,
			1,
			true,
		},
		{
			"wrong userID, user is not owner",
			1,
			2,
			false,
		},
		{
			"invalid userID, user is not owner",
			1,
			-1,
			false,
		},
		{
			"invalid taskID, user is not owner",
			-1,
			1,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, _ = testDB.Exec("TRUNCATE tasks, users RESTART IDENTITY CASCADE ")
			})

			ctx := context.Background()
			taskRepo := NewTaskRepository(*testDB)
			testCreateUser(t, *testDB)
			ctx = context.WithValue(ctx, contextkeys.UserID, int64(1))

			p := models.TaskPayload{
				Name:     "Lorem Ipsum",
				Priority: models.PriorityLow,
			}

			_ = storeTask(t, taskRepo, p, ctx)

			result, err := taskRepo.IsTaskOwner(tc.userID, tc.taskID)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if result != tc.expectedResult {
				t.Errorf("Wrong outcome, expected %t, got %t", tc.expectedResult, result)
			}
		})
	}
}
