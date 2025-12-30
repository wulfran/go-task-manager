package services

import (
	"context"
	"fmt"
	"strings"
	"task-manager/internal/helpers"
	"task-manager/internal/models"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type mockUserRepository struct {
}

func (m mockUserRepository) CreateUser(ctx context.Context, r models.CreateUserPayload) error {
	if r.Email == "fail-create@example.com" {
		return fmt.Errorf("failed to create a user")
	}
	return nil
}

func (m mockUserRepository) CheckIfEmailExists(email string) (bool, error) {
	switch email {
	case "example@test.com":
		return true, nil
	case "error@test.com":
		return false, fmt.Errorf("error while executing the query")
	default:
		return false, nil
	}
}

func (m mockUserRepository) GetUserData(p models.LoginPayload) (models.User, error) {
	switch p.Email {
	case "no-user-found@test.com":
		return models.User{}, fmt.Errorf("GetUserData: no entries found")
	case "test@example.com":
		return models.User{
			ID:        1,
			Name:      "Lorem Ipsum",
			Email:     "test@example.com",
			Password:  "loremIpsum",
			CreatedAt: nil,
		}, nil
	case "wrong-password@test.com":
		h, _ := helpers.HashPassword("loremIpsum")

		if !helpers.ValidatePassword(p.Password, h) {
			return models.User{}, fmt.Errorf("incorrect credentials")
		}
	default:
		return models.User{}, fmt.Errorf("unexpected error")
	}
	return models.User{}, nil
}

func TestUserService_RegisterUser(t *testing.T) {
	s := NewUserService(mockUserRepository{})
	var tests = []struct {
		name         string
		payload      models.CreateUserPayload
		expectsError bool
		wantedError  string
	}{
		{
			"valid payload, no errors",
			models.CreateUserPayload{
				Name:     "Lorem ipsum",
				Email:    "user@new.com",
				Password: "LoremIpsum",
			},
			false,
			"",
		},
		{
			"email exist, returns error",
			models.CreateUserPayload{
				Name:     "Lorem ipsum",
				Email:    "example@test.com",
				Password: "LoremIpsum",
			},
			true,
			"RegisterUser: email already in use",
		},
		{
			"password hash error",
			models.CreateUserPayload{
				Name:     "Lorem ipsum",
				Email:    "test@example.com",
				Password: strings.Repeat("a", 73),
			},
			true,
			"RegisterUser: failed to hash a password, bcrypt: password length exceeds 72 bytes",
		},
		{
			"error while checking email",
			models.CreateUserPayload{
				Name:     "Lorem ipsum",
				Email:    "error@test.com",
				Password: "LoremIpsum",
			},
			true,
			"RegisterUser: failed to check if email is unique: error while executing the query",
		},
		{
			"error while creating user",
			models.CreateUserPayload{
				Name:     "Lorem ipsum",
				Email:    "fail-create@example.com",
				Password: "LoremIpsum",
			},
			true,
			"RegisterUser: failed to run create user: failed to create a user",
		},
	}
	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			err := s.RegisterUser(context.Background(), i.payload)
			if i.expectsError && err == nil {
				t.Errorf("RegisterUser should return error but it did not")
			}

			if i.expectsError {
				if !strings.Contains(err.Error(), i.wantedError) {
					t.Errorf("expected error <%s> but got <%s>", i.wantedError, err)
				}
			}

			if !i.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	s := NewUserService(mockUserRepository{})
	var tests = []struct {
		name                    string
		payload                 models.LoginPayload
		expectedPayloadReturned models.User
		expectError             bool
		errorWanted             string
	}{
		{
			"correct login",
			models.LoginPayload{
				Email:    "test@example.com",
				Password: "loremIpsum",
			},
			models.User{
				ID:        1,
				Name:      "Lorem Ipsum",
				Email:     "test@example.com",
				Password:  "loremIpsum",
				CreatedAt: nil,
			},
			false,
			"",
		},
		{
			"no user found",
			models.LoginPayload{
				Email:    "no-user-found@test.com",
				Password: "DolorEt",
			},
			models.User{
				ID:        0,
				Name:      "",
				Email:     "",
				Password:  "",
				CreatedAt: nil,
			},
			true,
			"LoginUser: failed to get user data: GetUserData: no entries found",
		},
		{
			"incorrect password",
			models.LoginPayload{
				Email:    "wrong-password@test.com",
				Password: "DolorEt",
			},
			models.User{
				ID:        0,
				Name:      "",
				Email:     "",
				Password:  "",
				CreatedAt: nil,
			},
			true,
			"LoginUser: failed to get user data: incorrect credentials",
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			u, err := s.LoginUser(i.payload)
			if i.expectError && err == nil {
				t.Errorf("function is expected to return an error but it did not")
			}
			if !i.expectError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if !i.expectError {
				if diff := cmp.Diff(i.expectedPayloadReturned, u); diff != "" {
					t.Errorf("invalid data returned, (-want, +got)\n %s", diff)
				}
			} else {
				if !strings.Contains(err.Error(), i.errorWanted) {
					t.Errorf("expected error <%s> but got <%s>", i.errorWanted, err)
				}
			}
		})
	}
}

func TestUserService_CheckIfEmailExists(t *testing.T) {
	s := NewUserService(mockUserRepository{})

	var tests = []struct {
		name           string
		testedEmail    string
		expectedResult bool
	}{
		{
			"email exists, returns true",
			"example@test.com",
			true,
		},
		{
			"email does not exist, returns false",
			"test@example.com",
			false,
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			exists, _ := s.CheckIfEmailExists(i.testedEmail)

			if i.expectedResult != exists {
				t.Errorf("CheckIfEmailExists(%q) = %t, wanted %t ", i.testedEmail, i.expectedResult, exists)
			}
		})
	}
}
