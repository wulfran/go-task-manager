package services

import (
	"task-manager/internal/config"
	"task-manager/internal/models"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthService_CreateToken(t *testing.T) {
	c := config.JWTConfig{Secret: "secret-for-testing"}

	var tests = []struct {
		name         string
		inputUser    models.User
		expectsError bool
		errorWanted  string
	}{
		{
			"valid data, correct output",
			models.User{
				ID:    1,
				Name:  "Lorem Ipsum",
				Email: "lorem@ipsum.com",
			},
			false,
			"",
		},
		{
			"invalid data",
			models.User{},
			true,
			"invalid user data provided",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewAuthService(c)

			token, err := s.CreateToken(tc.inputUser)
			if tc.expectsError {
				if err == nil {
					t.Errorf("function was supposed to return an error but it did not")
				}

				if err != nil && tc.errorWanted != err.Error() {
					t.Errorf("invalid error message, expected <%s> but got <%s>", tc.errorWanted, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}

				if token == "" {
					t.Errorf("token should not be empty")
				}

				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, jwt.ErrSignatureInvalid
					}

					return []byte(c.Secret), nil
				})

				if err != nil {
					t.Errorf("failed to parse token: %s", err)
					return
				}

				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					t.Errorf("failed to read claims")
					return
				}

				uName := claims["username"].(string)
				if uName != tc.inputUser.Email {
					t.Errorf("invalid user name in claims, expected %s but got %s", tc.inputUser.Name, uName)
				}

				uID := claims["userId"].(float64)
				if int(uID) != tc.inputUser.ID {
					t.Errorf("invalid user name in claims, expected %d but got %f", tc.inputUser.ID, uID)
				}

				exp, ok := claims["exp"].(float64)
				if !ok {
					t.Errorf("expiration is missing from the claims or is invalid")
				}

				expTime := time.Unix(int64(exp), 0)
				now := time.Now()
				if expTime.Before(now) {
					t.Errorf("expiration is in the past, and it should not be: %v", expTime)
				}
			}
		})
	}
}
