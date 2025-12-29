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
	if email != "example@test.com" {
		return false, nil
	}

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
