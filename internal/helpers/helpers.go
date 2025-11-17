package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func SliceContains[T comparable](slice []T, value T) bool {
	for _, val := range slice {
		if val == value {
			return true
		}
	}
	return false
}

func GetQueryPath(n string) string {
	return fmt.Sprintf("./internal/db/queries/" + n)
}

func HashPassword(p string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 6)
	if err != nil {
		fmt.Println("Failed to hash password!")
	}

	return string(bytes)
}

func JsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	tmp := emailRegex.MatchString(addr.Address)
	return tmp
}
