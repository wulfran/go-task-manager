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

	return r
}
