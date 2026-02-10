package repository

import (
	"testing"
)

func TestNewTaskRepository(t *testing.T) {
	repository := New(*testDB)

	if repository.Ur == nil {
		t.Errorf("userRepository should not be nil")
	}

	if repository.Tr == nil {
		t.Errorf("taskRepository should not be nil")
	}
}
