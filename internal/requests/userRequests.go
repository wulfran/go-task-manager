package requests

import (
	"task-manager/internal/helpers"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r CreateUserRequest) Validate() ValidationResult {
	var res ValidationResult

	if len(r.Name) < 3 {
		res.SetFailed("Name has to be at least 3 characters long.")
	}

	if !helpers.IsValidEmail(r.Email) {
		res.SetFailed("Email invalid.")
	}

	if len(r.Password) < 5 {
		res.SetFailed("Password has to be at least 5 characters long.")
	}

	return res
}

type UpdateUserRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c Credentials) Validate() ValidationResult {
	var r ValidationResult
	r.Validated = true
	if len(c.Email) < 1 {
		r.SetFailed("Missing e-mail")
	}

	if len(c.Password) < 1 {
		r.SetFailed("Missing password")
	}

	if !helpers.IsValidEmail(c.Email) {
		r.SetFailed("E-mail invalid")
	}

	return r
}
