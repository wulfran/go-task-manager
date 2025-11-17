package repository

import "task-manager/internal/db"

type Repositories struct {
	Ur UserRepository
}

func New(d db.DB) Repositories {
	return Repositories{
		Ur: NewUserRepository(d),
	}
}
