package requests

import (
	"strings"
	"task-manager/internal/helpers"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r CreateUserRequest) Validate() ValidationResult {
	res := ValidationResult{
		Validated: true,
		Message:   "",
	}

	nl := len(strings.TrimSpace(r.Name))
	if nl < 3 {
		res.SetFailed("name has to be at least 3 characters long")
	}

	if !helpers.IsValidEmail(r.Email) {
		res.SetFailed("email invalid")
	}

	pl := len(strings.TrimSpace(r.Password))
	if pl < 5 {
		res.SetFailed("password has to be at least 5 characters long")
	}
	if pl > 72 {
		res.SetFailed("password too long")
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
	r := ValidationResult{
		Validated: true,
		Message:   "",
	}
	el := len(strings.TrimSpace(c.Email))
	if el < 1 {
		r.SetFailed("missing e-mail")
	}

	pl := len(strings.TrimSpace(c.Password))
	if pl < 1 || pl > 72 {
		r.SetFailed("password is missing or has invalid length")
	}

	if !helpers.IsValidEmail(c.Email) {
		r.SetFailed("e-mail invalid")
	}

	return r
}
