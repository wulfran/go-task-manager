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
	Login() func(w http.ResponseWriter, r *http.Request)
}

type usersController struct {
	us services.UserService
	as services.AuthService
}

func NewUsersController(us services.UserService, as services.AuthService) UsersController {
	return &usersController{
		us: us,
		as: as,
	}
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (uc usersController) Store(BodySizeLimit int64) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, BodySizeLimit)

		var request requests.CreateUserRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			helpers.JsonResponse(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("store user: request body too large: %v", err))
			return
		}

		v := request.Validate()
		if !v.Validated {
			helpers.JsonResponse(w, http.StatusUnprocessableEntity, fmt.Sprintf("store user validation failed: %s", v.Message))
			return
		}

		payload := models.CreateUserPayload{
			Name:     request.Name,
			Email:    request.Email,
			Password: request.Password,
		}

		if err := uc.us.RegisterUser(r.Context(), payload); err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("store user: failed to register: %v", err))
			return
		}

		helpers.JsonResponse(w, http.StatusOK, fmt.Sprintf("successfully created user"))
	}
}
func (uc usersController) Login() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.Credentials
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.JsonResponse(w, http.StatusBadRequest, fmt.Sprintf("login error, incorrect payload: %v", err))
			return
		}

		v := req.Validate()
		if !v.Validated {
			helpers.JsonResponse(w, http.StatusUnprocessableEntity, fmt.Sprintf("login: login request is invalid: %s", v.Message))
			return
		}

		p := models.LoginPayload{
			Email:    req.Email,
			Password: req.Password,
		}

		u, err := uc.us.LoginUser(p)
		if err != nil {
			helpers.JsonResponse(w, http.StatusUnauthorized, fmt.Sprintf("failed to authenticate, incorrect credentials"))
			return
		}
		fmt.Printf("%v", u)

		token, err := uc.as.CreateToken(u)
		if err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to generate JWT token, %v", err))
			return
		}

		helpers.JsonResponse(w, http.StatusOK, AuthResponse{Token: token})
		return
	}
}
