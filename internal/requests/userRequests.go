package requests

import (
	"fmt"
	"task-manager/internal/helpers"
)

type ValidationResult struct {
	Validated bool
	Message   string
}
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

	if len(r.Name) < 3 {
		res.Validated = false
		res.Message = fmt.Sprintf(res.Message + " Name has to be at least 3 characters long.")
	}

	if !helpers.IsValidEmail(r.Email) {
		res.Validated = false
		res.Message = fmt.Sprintf(res.Message + " Email invalid.")
	}

	if len(r.Password) < 5 {
		res.Validated = false
		res.Message = fmt.Sprintf(res.Message + " Password has to be at least 5 characters long.")
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
		r.Validated = false
		r.Message = fmt.Sprintf(r.Message + " Missing email")
	}

	if len(c.Password) < 1 {
		r.Validated = false
		r.Message = fmt.Sprintf(r.Message + " Missing password")
	}

	if !helpers.IsValidEmail(c.Email) {
		r.Validated = false
		r.Message = fmt.Sprintf(r.Message + " Email invalid.")
	}

	return r
}
