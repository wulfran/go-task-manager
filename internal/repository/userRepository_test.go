package repository

import (
	"context"
	"task-manager/internal/helpers"
	"task-manager/internal/models"
	"testing"
)

func testCreateUser(t *testing.T, ur UserRepository) {
	t.Helper()

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
		name         string
		payload      models.CreateUserPayload
		expectsError bool
		errorWanted  string
	}{
		{
			"valid payload",
			models.CreateUserPayload{
				Name:     "Lorem Ipsum",
				Email:    "lorem@ipsum.com",
				Password: hash,
			},
			false,
			"",
		},
		{
			"existing email, should error",
			models.CreateUserPayload{
				Name:     "Lorem Ipsum",
				Email:    "lorem@ipsum.com",
				Password: hash,
			},
			true,
			"CreateUser: failed to insert a new user: pq: duplicate key value violates unique constraint \"users_email_key\"",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := userRepo.CreateUser(ctx, tc.payload)

			if tc.expectsError {
				if err == nil {
					t.Error("createUser should return an error but it did not")
				}

				if tc.errorWanted != err.Error() {
					t.Errorf("wrong error returned, expected <%s> but got <%s>", tc.errorWanted, err)
				}
			}

			if !tc.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
		})
	}
}

func TestUserRepository_CheckIfEmailExists(t *testing.T) {
	userRepo := NewUserRepository(*testDB)
	testCreateUser(t, userRepo)
	var tests = []struct {
		name           string
		email          string
		expectedResult bool
	}{
		{
			"no-existing e-mail, no error",
			"john.doe@example.com",
			false,
		},
		{
			"existing e-mail, no error",
			"lorem@ipsum.com",
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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
	testCreateUser(t, userRepo)
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
