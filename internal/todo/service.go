package todo

import (
	"context"
	"errors"
)

var (
	ErrTodoNotFound    = errors.New("todo not found")
	ErrInvalidTitle    = errors.New("title is required")
	ErrInvalidPriority = errors.New("invalid priority")
)

type Repository interface {
	Create(ctx context.Context, todo *Todo) error
	GetById(ctx context.Context, id int64) (*Todo, error)
	GetAll(ctx context.Context, priority *Priority) ([]Todo, error)
	Update(ctx context.Context, id int64, todo *Todo) error
}

type TodoService struct {
	repo Repository
}

func NewTodoService(repo Repository) *TodoService {
	return &TodoService{repo}
}

func (s *TodoService) Create(ctx context.Context, req *CreateTodoRequest) (*Todo, error) {
	if req.Title == "" {
		return nil, ErrInvalidTitle
	}

	if req.Priority == "" {
		req.SetDefaults()
	} else if !isValidPriority(req.Priority) {
		return nil, ErrInvalidPriority
	}

	todo := &Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Deadline:    req.Deadline,
		Completed:   false,
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) GetById(ctx context.Context, id int64) (*Todo, error) {
	todo, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	} else if todo == nil {
		return nil, ErrTodoNotFound
	}

	return todo, nil
}

func (s *TodoService) GetAll(ctx context.Context, priority *Priority) ([]Todo, error) {
	if priority != nil && !isValidPriority(*priority) {
		return nil, ErrInvalidPriority
	}
	return s.repo.GetAll(ctx, priority)
}

func (s *TodoService) Update(ctx context.Context, id int64, req *UpdateTodoRequest) (*Todo, error) {
	todo, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, ErrTodoNotFound
	}

	if req.Title == "" {
		return nil, ErrInvalidTitle
	} else {
		todo.Title = req.Title
	}

	if req.Description != nil {
		todo.Description = req.Description
	}

	if req.Priority != nil {
		if !isValidPriority(*req.Priority) {
			return nil, ErrInvalidPriority
		}
		todo.Priority = *req.Priority
	}

	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	if err := s.repo.Update(ctx, id, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func isValidPriority(p Priority) bool {
	switch p {
	case Low, Medium, High:
		return true
	default:
		return false
	}
}
