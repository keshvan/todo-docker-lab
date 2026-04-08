package todo_test

import (
	"context"
	"testing"
	"time"

	todopkg "github.com/keshvan/todo-docker-lab/internal/todo"
	"github.com/keshvan/todo-docker-lab/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTodoService_Create(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	svc := todopkg.NewTodoService(mockRepo)

	desc := "description"
	deadline := time.Now()

	req := &todopkg.CreateTodoRequest{
		Title:       "Test Task",
		Description: &desc,
		Priority:    todopkg.High,
		Deadline:    &deadline,
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(todo *todopkg.Todo) bool {
		return todo.Title == req.Title && todo.Priority == req.Priority
	})).Return(nil)

	todo, err := svc.Create(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, req.Title, todo.Title)
	require.Equal(t, req.Priority, todo.Priority)
	require.Equal(t, &desc, todo.Description)

	mockRepo.AssertExpectations(t)
}

func TestTodoService_Create_InvalidTitle(t *testing.T) {
	svc := todopkg.NewTodoService(nil)

	req := &todopkg.CreateTodoRequest{
		Title: "",
	}

	todo, err := svc.Create(context.Background(), req)

	require.Nil(t, todo)
	require.ErrorIs(t, err, todopkg.ErrInvalidTitle)
}

func TestTodoService_Create_InvalidPriority(t *testing.T) {
	svc := todopkg.NewTodoService(nil)

	req := &todopkg.CreateTodoRequest{
		Title:    "Task",
		Priority: "invalid",
	}

	todo, err := svc.Create(context.Background(), req)

	require.Nil(t, todo)
	require.ErrorIs(t, err, todopkg.ErrInvalidPriority)
}

func TestTodoService_GetById(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	svc := todopkg.NewTodoService(mockRepo)

	todoFromRepo := &todopkg.Todo{ID: 1, Title: "Task"}

	mockRepo.On("GetById", mock.Anything, int64(1)).Return(todoFromRepo, nil)

	todo, err := svc.GetById(context.Background(), 1)

	require.NoError(t, err)
	require.Equal(t, int64(1), todo.ID)

	mockRepo.AssertExpectations(t)
}

func TestTodoService_GetById_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	svc := todopkg.NewTodoService(mockRepo)

	mockRepo.On("GetById", mock.Anything, int64(1)).Return(nil, nil)

	todo, err := svc.GetById(context.Background(), 1)

	require.Nil(t, todo)
	require.ErrorIs(t, err, todopkg.ErrTodoNotFound)
}

func TestTodoService_GetAll(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	svc := todopkg.NewTodoService(mockRepo)

	todosFromRepo := []todopkg.Todo{
		{ID: 1, Title: "Task1", Priority: todopkg.High},
		{ID: 2, Title: "Task2", Priority: todopkg.Low},
	}

	priority := todopkg.High

	mockRepo.On("GetAll", mock.Anything, &priority).Return(todosFromRepo, nil)

	todos, err := svc.GetAll(context.Background(), &priority)

	require.NoError(t, err)
	require.Len(t, todos, 2)
	mockRepo.AssertExpectations(t)
}

func TestTodoService_GetAll_InvalidPriority(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	svc := todopkg.NewTodoService(mockRepo)

	invalid := todopkg.Priority("invalid")
	todos, err := svc.GetAll(context.Background(), &invalid)

	require.Nil(t, todos)
	require.ErrorIs(t, err, todopkg.ErrInvalidPriority)
}

func TestTodoService_Update(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	svc := todopkg.NewTodoService(mockRepo)

	desc := "new desc"
	priority := todopkg.Medium
	completed := true

	existingTodo := &todopkg.Todo{
		ID:       1,
		Title:    "Old Title",
		Priority: todopkg.High,
	}

	req := &todopkg.UpdateTodoRequest{
		Title:       "Updated Title",
		Description: &desc,
		Priority:    &priority,
		Completed:   &completed,
	}

	mockRepo.On("GetById", mock.Anything, int64(1)).Return(existingTodo, nil)
	mockRepo.On("Update", mock.Anything, int64(1), existingTodo).Return(nil)

	todo, err := svc.Update(context.Background(), 1, req)

	require.NoError(t, err)
	require.Equal(t, "Updated Title", todo.Title)
	require.Equal(t, &desc, todo.Description)
	require.Equal(t, priority, todo.Priority)
	require.Equal(t, completed, todo.Completed)

	mockRepo.AssertExpectations(t)
}

func TestTodoService_Update_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	svc := todopkg.NewTodoService(mockRepo)

	req := &todopkg.UpdateTodoRequest{
		Title: "Updated",
	}

	mockRepo.On("GetById", mock.Anything, int64(1)).Return(nil, nil)

	todo, err := svc.Update(context.Background(), 1, req)

	require.Nil(t, todo)
	require.ErrorIs(t, err, todopkg.ErrTodoNotFound)
}
