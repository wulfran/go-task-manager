package models

import "time"

type User struct {
	ID        int        `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"password" db:"password"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}

type UserList struct {
	Users []User `json:"users"`
}

type CreateUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"Password"`
}
