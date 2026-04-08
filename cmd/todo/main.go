package main

import (
	"log"

	"github.com/keshvan/todo-docker-lab/internal/app"
	"github.com/keshvan/todo-docker-lab/internal/config"
)

func main() {
	cfg := config.MustLoad()

	a, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	a.Run(cfg.ServerPort)
}
