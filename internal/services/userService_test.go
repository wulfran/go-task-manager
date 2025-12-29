package services

import (
	"context"
	"task-manager/internal/models"
	"testing"
)

type mockUserRepository struct {
}

func (m mockUserRepository) CreateUser(ctx context.Context, r models.CreateUserPayload) error {
	return nil
}

func (m mockUserRepository) CheckIfEmailExists(email string) (bool, error) {
	return true, nil
}

func (m mockUserRepository) GetUserData(p models.LoginPayload) (models.User, error) {
	return models.User{}, nil
}

func TestUserService_CheckIfEmailExists(t *testing.T) {
	s := NewUserService(mockUserRepository{})

	var tests = []struct {
		name           string
		testedEmail    string
		expectedResult bool
		expectsError   bool
		wantedError    string
	}{
		{
			"valid email, returns true",
			"example@test.com",
			true,
			false,
			"",
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			exists, err := s.CheckIfEmailExists(i.testedEmail)

			if !i.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if i.expectsError && err == nil {
				t.Errorf("function was expected to return error but it did not")
			}

			if i.expectedResult != exists {
				t.Errorf("function supposed to return %t but got %t instead", i.expectedResult, exists)
			}
		})
	}
}
