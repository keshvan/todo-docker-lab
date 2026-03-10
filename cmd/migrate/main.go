package main

import (
	"errors"
	"log"
	"todo-api/internal/config"
	"todo-api/internal/todo"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	cfg := config.MustLoad()

	if err := todo.RunMigrations(cfg.DatabaseURL); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no new migrations")
		} else {
			log.Fatalf("erorr while applying migrations: %v", err)
		}
	}
}
