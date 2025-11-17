package routes

import (
	"fmt"
	"net/http"
	"task-manager/internal/controllers"
	"task-manager/internal/db"
	"task-manager/internal/helpers"
	"task-manager/internal/repository"
	"task-manager/internal/services"

	"github.com/go-chi/chi/v5"
	chimiddlware "github.com/go-chi/chi/v5/middleware"
)

const (
	bodySizeLimit = 50 << 20
)

func CreateServer(d db.DB) *chi.Mux {
	fmt.Println("Setting up repository, service and controller")
	ur := repository.NewUserRepository(d)
	us := services.NewUserService(ur)
	uc := controllers.NewUsersController(us)

	r := chi.NewRouter()

	r.Use(chimiddlware.RequestID)
	r.Use(chimiddlware.Recoverer)
	r.Use(chimiddlware.Logger)
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		helpers.JsonResponse(w, 405, fmt.Sprintf("method not allowed"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	r.Post("/register", uc.Store(bodySizeLimit))

	return r
}
