package services

import "task-manager/internal/repository"

type Services struct {
	Us UserService
	As AuthService
	Ts TaskService
}

func New(r repository.Repositories) Services {
	return Services{
		Us: NewUserService(r.Ur),
		As: NewAuthService(),
		Ts: NewTaskService(r.Tr),
	}
}
