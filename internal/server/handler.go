package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"todo-api/internal/todo"

	"github.com/go-chi/chi/v5"
)

type Service interface {
	Create(ctx context.Context, req *todo.CreateTodoRequest) (*todo.Todo, error)
	GetById(ctx context.Context, id int64) (*todo.Todo, error)
	GetAll(ctx context.Context, priority *todo.Priority) ([]todo.Todo, error)
	Update(ctx context.Context, id int64, todo *todo.UpdateTodoRequest) (*todo.Todo, error)
}

type TodoHandler struct {
	service Service
}

func NewTodoHandler(service Service) *TodoHandler {
	return &TodoHandler{service}
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req todo.CreateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	todo, err := h.service.Create(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	todo, err := h.service.GetById(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var priority *todo.Priority

	p := r.URL.Query().Get("priority")
	if p != "" {
		val := todo.Priority(p)
		priority = &val
	}

	todos, err := h.service.GetAll(r.Context(), priority)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req todo.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	todo, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var code int
	var msg string

	switch {
	case errors.Is(err, todo.ErrTodoNotFound):
		code = http.StatusNotFound
		msg = err.Error()
	case errors.Is(err, todo.ErrInvalidPriority),
		errors.Is(err, todo.ErrInvalidTitle):
		code = http.StatusBadRequest
		msg = err.Error()
	default:
		code = http.StatusInternalServerError
		msg = "internal server error"
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
