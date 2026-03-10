package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	chi *chi.Mux
}

func NewRouter(todoHandler Handler) *Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/todos", func(r chi.Router) {
			r.Post("/", todoHandler.Create)
			r.Get("/", todoHandler.GetAll)
			r.Get("/{id}", todoHandler.GetByID)
			r.Patch("/{id}", todoHandler.Update)
		})
	})

	return &Router{r}
}

func (r *Router) Handler() http.Handler {
	return r.chi
}
