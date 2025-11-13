package main

import (
	"log"
	"task-manager/internal/db"
)

func main() {
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
}
