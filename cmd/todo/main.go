package main

import (
	"log"
	"todo-api/internal/app"
	"todo-api/internal/config"
)

func main() {
	cfg := config.MustLoad()

	a, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	a.Run(cfg.ServerPort)
}
