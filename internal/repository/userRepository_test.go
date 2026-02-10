package repository

import (
	"context"
	"task-manager/internal/db"
	"task-manager/internal/helpers"
	"task-manager/internal/models"
	"testing"
)

func testCreateUser(t *testing.T, db db.DB) {
	t.Helper()

	ur := NewUserRepository(db)

	ctx := context.Background()
	hash, _ := helpers.HashPassword("secretPassword")

	userData := models.CreateUserPayload{
		Name:     "Lorem Ipsum",
		Email:    "lorem@ipsum.com",
		Password: hash,
	}
	_ = ur.CreateUser(ctx, userData)
}

func TestUserRepository_CreateUser(t *testing.T) {
	userRepo := NewUserRepository(*testDB)
	hash, _ := helpers.HashPassword("secretPassword")

	var tests = []struct {
		name                    string
		payload                 models.CreateUserPayload
		shouldCreateDefaultUser bool
		expectsError            bool
		errorWanted             string
	}{
		{
			"valid payload",
			models.CreateUserPayload{
				Name:     "Lorem Ipsum",
				Email:    "lorem@ipsum.com",
				Password: hash,
			},
			false,
			false,
			"",
		},
		{
			"existing email should error",
			models.CreateUserPayload{
				Name:     "Lorem Ipsum",
				Email:    "lorem@ipsum.com",
				Password: hash,
			},
			true,
			true,
			"CreateUser: failed to insert a new user: pq: duplicate key value violates unique constraint \"users_email_key\"",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, _ = testDB.Exec("TRUNCATE users RESTART IDENTITY CASCADE ")
			})
			ctx := context.Background()
			if tc.shouldCreateDefaultUser {
				testCreateUser(t, *testDB)
			}
			err := userRepo.CreateUser(ctx, tc.payload)
			if tc.expectsError {
				if err == nil {
					t.Error("function should return an error but it did not")
				}

				if err != nil && tc.errorWanted != err.Error() {
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

func TestUserRepository_CheckIfEmailExists(t *testing.T) {
	userRepo := NewUserRepository(*testDB)
	testCreateUser(t, *testDB)
	var tests = []struct {
		name                    string
		email                   string
		shouldCreateDefaultUser bool
		expectedResult          bool
	}{
		{
			"no-existing e-mail, no error",
			"john.doe@example.com",
			false,
			false,
		},
		{
			"existing e-mail, no error",
			"lorem@ipsum.com",
			true,
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, _ = testDB.Exec("TRUNCATE users RESTART IDENTITY CASCADE")
			})
			if tc.shouldCreateDefaultUser {
				testCreateUser(t, *testDB)
			}

			result, err := userRepo.CheckIfEmailExists(tc.email)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if result != tc.expectedResult {
				t.Errorf("wrong result, expected %t but got %t", tc.expectedResult, result)
			}
		})
	}
}

func TestUserRepository_GetUserData(t *testing.T) {
	userRepo := NewUserRepository(*testDB)
	testCreateUser(t, *testDB)
	hash, _ := helpers.HashPassword("secretPassword")

	var tests = []struct {
		name           string
		payload        models.LoginPayload
		expectedResult models.User
		expectsError   bool
		errorWanted    string
	}{
		{
			"valid payload",
			models.LoginPayload{
				Email:    "lorem@ipsum.com",
				Password: "secretPassword",
			},
			models.User{
				ID:    1,
				Name:  "Lorem Ipsum",
				Email: "lorem@ipsum.com",
			},
			false,
			"",
		},
		{
			"invalid credentials",
			models.LoginPayload{
				Email:    "lorem@ipsum.com",
				Password: "loremIpsum",
			},
			models.User{},
			true,
			"incorrect credentials",
		},
		{
			"password missing",
			models.LoginPayload{
				Email: "lorem@ipsum.com",
			},
			models.User{},
			true,
			"incorrect credentials",
		},
		{
			"email missing",
			models.LoginPayload{
				Password: "secretPassword",
			},
			models.User{},
			true,
			"GetUserData: no entries found",
		},
		{
			"hashed password provided",
			models.LoginPayload{
				Email:    "lorem@ipsum.com",
				Password: hash,
			},
			models.User{},
			true,
			"incorrect credentials",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user, err := userRepo.GetUserData(tc.payload)
			if tc.expectsError {
				if err == nil {
					t.Error("getUserData was supposed to return an error but it did not")
				} else if tc.errorWanted != err.Error() {
					t.Errorf("wrong error returned, expected <%s> but got <%s>", tc.errorWanted, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
				if user.ID != tc.expectedResult.ID || user.Email != tc.expectedResult.Email || user.Name != tc.expectedResult.Name {
					t.Errorf("wrong user returned: %v", user)
				}
			}
		})
	}
}
