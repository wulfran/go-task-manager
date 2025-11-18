package server

import (
	"fmt"
	"net/http"
	"task-manager/internal/helpers"

	"github.com/go-chi/chi/v5"
	chimiddlware "github.com/go-chi/chi/v5/middleware"
)

const (
	bodySizeLimit = 50 << 20
)

func (s Server) CreateServer() *chi.Mux {
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
	r.Post("/register", s.C.Uc.Store(bodySizeLimit))
	r.Post("/login", s.C.Uc.Login())

	r.Group(func(r chi.Router) {
		r.Use(s.Authenticate)
		r.Get("/marco", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("polo!"))
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", s.C.Tc.Index())
			r.Get("/{task_id}", s.C.Tc.Show())
			r.Post("/", s.C.Tc.Store(bodySizeLimit))
			r.Patch("/{task_id}", s.C.Tc.Update(bodySizeLimit))
			r.Delete("/{task_id}", s.C.Tc.Delete())
		})
	})

	return r
}
