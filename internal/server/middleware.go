package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"task-manager/internal/contextkeys"
	"task-manager/internal/helpers"

	"github.com/golang-jwt/jwt/v5"
)

func (s Server) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtSecret := []byte(s.Cfg.JWT.Secret)
		tString := r.Header.Get("Authorization")
		if tString == "" {
			helpers.JsonResponse(w, http.StatusUnauthorized, fmt.Sprintf("token is missing!"))
			return
		}
		if len(tString) > 7 && strings.ToUpper(tString[0:6]) == "BEARER" {
			tString = tString[7:]
		}

		t, err := jwt.Parse(tString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("authenticate: invalid sign method")
			}
			return jwtSecret, nil
		})

		if err != nil || !t.Valid {
			helpers.JsonResponse(w, http.StatusUnauthorized, fmt.Sprintf("invalid token"))
			return
		}
		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			helpers.JsonResponse(w, http.StatusUnauthorized, fmt.Sprintf("invalid token"))
			return
		}
		email := claims["username"].(string)
		fID := claims["userId"].(float64)
		uID := int64(fID)

		exists, err := s.S.Us.CheckIfEmailExists(email)
		if err != nil || !exists {
			helpers.JsonResponse(w, http.StatusUnauthorized, fmt.Sprintf("unable to verify user within the token"))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextkeys.UserID, uID)))
	})
}
