package services

import (
	"context"
	"fmt"
	"task-manager/internal/models"
	"task-manager/internal/repository"
)

type TaskService interface {
	StoreTask(ctx context.Context, p models.TaskPayload) error
	UpdateTask(ctx context.Context, p models.UpdateTask) (models.Task, error)
	ShowTask(id int) (models.Task, error)
	GetTasksList(uID int64) (models.TasksList, error)
	DeleteTask(id int, uID int64) error
	IsTaskOwner(uID int64, id int) (bool, error)
}

type taskService struct {
	r repository.TaskRepository
}

func NewTaskService(r repository.TaskRepository) TaskService {
	return &taskService{r: r}
}

func (s taskService) GetTasksList(uID int64) (models.TasksList, error) {
	if uID < 1 {
		return models.TasksList{}, fmt.Errorf("GetTasksList: invalid user")
	}

	l, err := s.r.Index(uID)
	if err != nil {
		return models.TasksList{}, fmt.Errorf("GetTasksList: failed to get data: %v", err)
	}

	return l, nil
}
func (s taskService) StoreTask(ctx context.Context, p models.TaskPayload) error {
	if err := s.r.Store(ctx, p); err != nil {
		return fmt.Errorf("storeTask: error while storing the data: %v", err)
	}

	return nil
}
func (s taskService) UpdateTask(ctx context.Context, p models.UpdateTask) (models.Task, error) {
	t, err := s.r.Update(ctx, p)
	if err != nil {
		return models.Task{}, fmt.Errorf("UpdateTask: %v", err)
	}
	return t, nil
}
func (s taskService) ShowTask(id int) (models.Task, error) {
	t, err := s.r.Show(id)
	if err != nil {
		return models.Task{}, fmt.Errorf("ShowTask: %v", err)
	}
	return t, nil
}
func (s taskService) DeleteTask(id int, uID int64) error {
	isOwner, err := s.IsTaskOwner(uID, id)
	if err != nil {
		return fmt.Errorf("deleteTask: %s", err)
	}
	if !isOwner {
		return fmt.Errorf("deleteTask: you are not authorized to execute this action")
	}

	if err := s.r.Delete(id); err != nil {
		return fmt.Errorf("DeleteTask: %s", err)
	}

	return nil
}
func (s taskService) IsTaskOwner(uID int64, id int) (bool, error) {
	isOwner, err := s.r.IsTaskOwner(uID, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is a task owner: %v", err)
	}
	return isOwner, nil
}
