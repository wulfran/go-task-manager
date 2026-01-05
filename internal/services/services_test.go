package services

import (
	"task-manager/internal/config"
	"task-manager/internal/repository"
	"testing"
)

func TestService_New(t *testing.T) {
	mockUserRepository := mockUserRepository{}
	mockTaskRepository := mockTaskRepository{}

	r := repository.Repositories{
		Ur: mockUserRepository,
		Tr: mockTaskRepository,
	}
	c := config.JWTConfig{Secret: "example-secret-for-testing"}

	s := New(r, c)

	if s.Us == nil {
		t.Errorf("userService should not be nil")
	}

	if s.Ts == nil {
		t.Errorf("taskService should not be nil")
	}

	if s.As == nil {
		t.Errorf("authService should not be nil")
	}
}
