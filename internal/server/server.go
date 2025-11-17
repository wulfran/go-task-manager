package server

import (
	"fmt"
	"log"
	"task-manager/internal/controllers"
	"task-manager/internal/db"
	"task-manager/internal/repository"
	"task-manager/internal/services"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	D db.DB
	C controllers.Controllers
	S services.Services
	R repository.Repositories
	H *chi.Mux
}

func New() *Server {
	// initialize database
	d, err := db.InitDb()
	if err != nil {
		log.Fatalf("unabled to initialize database: %v", err)
	}

	// test connection
	if err := d.Ping(); err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}

	// run migrations
	if err := d.RunMigrations(); err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Setting up repository, service and controller")
	r := repository.New(*d)
	svs := services.New(r)
	c := controllers.New(svs)

	s := &Server{
		D: *d,
		C: c,
		S: svs,
		R: r,
	}

	s.H = s.CreateServer()

	return s
}
