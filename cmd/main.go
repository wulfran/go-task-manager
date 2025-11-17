package main

import (
	"fmt"
	"net/http"
	"task-manager/internal/server"
)

func main() {
	s := server.New()

	fmt.Println("Launching the server")
	_ = http.ListenAndServe(":8000", s.H)
}
