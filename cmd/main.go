package main

import (
	"fmt"
	"net/http"
	"task-manager/internal/config"
	"task-manager/internal/server"
)

func main() {
	cfg := config.Load()

	s := server.New(cfg)

	fmt.Println("Launching the server")
	_ = http.ListenAndServe(":8000", s.H)
}
