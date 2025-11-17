package services

import "task-manager/internal/repository"

type Services struct {
	Us UserService
}

func New(r repository.Repositories) Services {
	return Services{
		Us: NewUserService(r.Ur),
	}
}
