package service

import (
	"context"
	"errors"
	"tasklist/internal/models"
	"tasklist/internal/repository"
)

type Task interface {
	Create(ctx context.Context, userId int, task models.TaskRequest) error
	Complete(ctx context.Context, taskId int) error
	List(ctx context.Context, userId int) ([]models.Task, error)
	GetByID(ctx context.Context, taskId, userId int) (*models.Task, error)
	Update(ctx context.Context, taskId, userId int, task models.TaskRequest) error
	Delete(ctx context.Context, id int) error
}
type TaskService struct {
	repo repository.Task
}

func NewTaskService(r repository.Task) *TaskService {
	return &TaskService{repo: r}
}

func (s *TaskService) Create(ctx context.Context, userId int, req models.TaskRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if err := s.repo.Create(ctx, userId, req); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) List(ctx context.Context, userId int) ([]models.Task, error) {
	return s.repo.List(ctx, userId)
}
func (s *TaskService) Complete(ctx context.Context, taskId int) error {
	return s.repo.Complete(ctx, taskId)
}

func (s *TaskService) GetByID(ctx context.Context, taskId, userId int) (*models.Task, error) {
	return s.repo.GetByID(ctx, taskId, userId)
}

func (s *TaskService) Update(ctx context.Context, taskId, userId int, task models.TaskRequest) error {
	return s.repo.Update(ctx, taskId, userId, task)
}

func (s *TaskService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
