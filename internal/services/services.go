package services

import (
	"task-manager/internal/config"
	"task-manager/internal/repository"
)

type Services struct {
	Us UserService
	As AuthService
	Ts TaskService
}

func New(r repository.Repositories, cfg config.JWTConfig) Services {
	return Services{
		Us: NewUserService(r.Ur),
		As: NewAuthService(cfg),
		Ts: NewTaskService(r.Tr),
	}
}
