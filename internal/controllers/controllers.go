package controllers

import "task-manager/internal/services"

type Controllers struct {
	Uc UsersController
	Tc TasksController
}

func New(s services.Services) Controllers {
	return Controllers{
		Uc: NewUsersController(s.Us, s.As),
		Tc: NewTasksController(s.Ts),
	}
}
