package requests

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateUserRequest_Validate(t *testing.T) {
	var tests = []struct {
		name           string
		credentials    CreateUserRequest
		expectedResult ValidationResult
	}{
		{
			"valid data",
			CreateUserRequest{
				Name:     "Lorem Ipsum",
				Email:    "lorem@ipsum.com",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: true,
				Message:   "",
			},
		},
		{
			"name empty string",
			CreateUserRequest{
				Name:     "",
				Email:    "lorem@ipsum.com",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "name has to be at least 3 characters long",
			},
		},
		{
			"missing name",
			CreateUserRequest{
				Email:    "lorem@ipsum.com",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "name has to be at least 3 characters long",
			},
		},
		{
			"name just whitespaces",
			CreateUserRequest{
				Name:     "   ",
				Email:    "lorem@ipsum.com",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "name has to be at least 3 characters long",
			},
		},
		{
			"invalid email",
			CreateUserRequest{
				Name:     "Lorem Ipsum",
				Email:    "loremipsum",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "email invalid",
			},
		},
		{
			"email missing",
			CreateUserRequest{
				Name:     "Lorem Ipsum",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "email invalid",
			},
		},
		{
			"email empty string",
			CreateUserRequest{
				Name:     "Lorem Ipsum",
				Email:    "",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "email invalid",
			},
		},
		{
			"email just whitespaces",
			CreateUserRequest{
				Name:     "Lorem Ipsum",
				Email:    "   ",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "email invalid",
			},
		},
		{
			"password missing",
			CreateUserRequest{
				Name:  "Lorem Ipsum",
				Email: "lorem@ipsum.com",
			},
			ValidationResult{
				Validated: false,
				Message:   "password has to be at least 5 characters long",
			},
		},
		{
			"password empty string",
			CreateUserRequest{
				Name:     "Lorem Ipsum",
				Email:    "lorem@ipsum.com",
				Password: "",
			},
			ValidationResult{
				Validated: false,
				Message:   "password has to be at least 5 characters long",
			},
		},
		{
			"password just whitespaces",
			CreateUserRequest{
				Name:     "Lorem Ipsum",
				Email:    "lorem@ipsum.com",
				Password: "     ",
			},
			ValidationResult{
				Validated: false,
				Message:   "password has to be at least 5 characters long",
			},
		},
		{
			"password too short",
			CreateUserRequest{
				Name:     "Lorem Ipsum",
				Email:    "lorem@ipsum.com",
				Password: "lore",
			},
			ValidationResult{
				Validated: false,
				Message:   "password has to be at least 5 characters long",
			},
		},
		{
			"password too long",
			CreateUserRequest{
				Name:     "Lorem Ipsum",
				Email:    "lorem@ipsum.com",
				Password: strings.Repeat("a", 73),
			},
			ValidationResult{
				Validated: false,
				Message:   "password too long",
			},
		},
		{
			"missing credentials entirely",
			CreateUserRequest{},
			ValidationResult{
				Validated: false,
				Message:   "name has to be at least 3 characters long, email invalid, password has to be at least 5 characters long",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.credentials.Validate()

			if diff := cmp.Diff(tc.expectedResult, res); diff != "" {
				t.Errorf("invalid validation result <-want, +got>\n%s", diff)
			}
		})
	}
}

func TestCredentials_Validate(t *testing.T) {
	var tests = []struct {
		name           string
		credentials    Credentials
		expectedResult ValidationResult
	}{
		{
			"valid credentials",
			Credentials{
				Email:    "lorem@ipsum.com",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: true,
				Message:   "",
			},
		},
		{
			"email missing",
			Credentials{
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "missing e-mail, e-mail invalid",
			},
		},
		{
			"email empty string",
			Credentials{
				Email:    "",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "missing e-mail, e-mail invalid",
			},
		},
		{
			"email just whitespaces",
			Credentials{
				Email:    "   ",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "missing e-mail, e-mail invalid",
			},
		},
		{
			"email present but invalid",
			Credentials{
				Email:    "loremIpsum",
				Password: "l0r3mIpsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "e-mail invalid",
			},
		},
		{
			"password missing",
			Credentials{
				Email: "lorem@ipsum.com",
			},
			ValidationResult{
				Validated: false,
				Message:   "password is missing or has invalid length",
			},
		},
		{
			"password empty string",
			Credentials{
				Email:    "lorem@ipsum.com",
				Password: "",
			},
			ValidationResult{
				Validated: false,
				Message:   "password is missing or has invalid length",
			},
		},
		{
			"password just whitespaces",
			Credentials{
				Email:    "lorem@ipsum.com",
				Password: "   ",
			},
			ValidationResult{
				Validated: false,
				Message:   "password is missing or has invalid length",
			},
		},
		{
			"password too long",
			Credentials{
				Email:    "lorem@ipsum.com",
				Password: strings.Repeat("a", 73),
			},
			ValidationResult{
				Validated: false,
				Message:   "password is missing or has invalid length",
			},
		},
		{
			"missing credentials entirely",
			Credentials{},
			ValidationResult{
				Validated: false,
				Message:   "missing e-mail, password is missing or has invalid length, e-mail invalid",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.credentials.Validate()

			if diff := cmp.Diff(tc.expectedResult, res); diff != "" {
				t.Errorf("invalid validation result <-want, +got>\n%s", diff)
			}
		})
	}
}
