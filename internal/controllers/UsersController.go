package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"task-manager/internal/helpers"
	"task-manager/internal/models"
	"task-manager/internal/requests"
	"task-manager/internal/services"
)

type UsersController interface {
	Store(BodySizeLimit int64) func(http.ResponseWriter, *http.Request)
}

type usersController struct {
	s services.UserService
}

func NewUsersController(s services.UserService) UsersController {
	return &usersController{s: s}
}

func (uc usersController) Store(BodySizeLimit int64) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, BodySizeLimit)

		var request requests.CreateUserRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			helpers.JsonResponse(w, 413, fmt.Sprintf("store user: Request Body Too Large"))
			return
		}

		v := request.Validate()
		if !v.Validated {
			helpers.JsonResponse(w, 422, fmt.Sprintf("store user validation failed: %s", v.Message))
			return
		}

		payload := models.CreateUserPayload{
			Name:     request.Name,
			Email:    request.Email,
			Password: request.Password,
		}

		if err := uc.s.RegisterUser(r.Context(), payload); err != nil {
			helpers.JsonResponse(w, 500, fmt.Sprintf("store user: failed to register: %v", err))
			return
		}

		helpers.JsonResponse(w, 200, fmt.Sprintf("successfully created user"))
	}
}
