package controllers

import "task-manager/internal/services"

type Controllers struct {
	Uc UsersController
}

func New(s services.Services) Controllers {
	return Controllers{
		Uc: NewUsersController(s.Us, s.As),
	}
}
