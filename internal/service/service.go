package service

import (
	"tasklist/internal/repository"
	"tasklist/pkg/config"
)

type Service struct {
	AuthService Auth
	TaskService Task
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		AuthService: NewAuthService(repo.UserRepo, cfg),
		TaskService: NewTaskService(repo.TaskRepo),
	}
}
