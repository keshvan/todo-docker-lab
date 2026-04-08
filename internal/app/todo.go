package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/keshvan/todo-docker-lab/internal/config"
	"github.com/keshvan/todo-docker-lab/internal/server"
	"github.com/keshvan/todo-docker-lab/internal/todo"
)

type App struct {
	Server *server.Server
	DB     *todo.Database
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := todo.New(cfg)
	if err != nil {
		return nil, err
	}

	repo := todo.NewTodoRepository(db.Pool)
	service := todo.NewTodoService(repo)
	handler := server.NewTodoHandler(service)
	router := server.NewRouter(handler)

	srv := server.NewServer(cfg.ServerPort, router.Handler())

	return &App{Server: srv, DB: db}, nil
}

func (a *App) Run(port string) {
	a.Server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}

	a.DB.Close()
	log.Println("Server exited properly")
}
