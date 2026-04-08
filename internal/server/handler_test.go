package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/keshvan/todo-docker-lab/internal/todo"
	"github.com/keshvan/todo-docker-lab/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func withURLParam(r *http.Request, key, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))
}

func TestTodoHandler_Create(t *testing.T) {
	svc := mocks.NewMockService(t)
	handler := NewTodoHandler(svc)

	reqBody := []byte(`{"title":"Test Task","priority":"high"}`)
	expectedReq := &todo.CreateTodoRequest{
		Title:    "Test Task",
		Priority: todo.High,
	}
	expectedTodo := &todo.Todo{
		ID:       1,
		Title:    "Test Task",
		Priority: todo.High,
	}

	svc.On("Create", mock.Anything, expectedReq).Return(expectedTodo, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var got todo.Todo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, *expectedTodo, got)

	svc.AssertExpectations(t)
}

func TestTodoHandler_Create_InvalidBody(t *testing.T) {
	handler := NewTodoHandler(mocks.NewMockService(t))

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid body")
}

func TestTodoHandler_Create_ServiceError(t *testing.T) {
	svc := mocks.NewMockService(t)
	handler := NewTodoHandler(svc)

	expectedReq := &todo.CreateTodoRequest{
		Title: "Task",
	}

	svc.On("Create", mock.Anything, expectedReq).Return(nil, todo.ErrInvalidTitle).Once()

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{"title":"Task"}`))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"error":"title is required"}`, rec.Body.String())

	svc.AssertExpectations(t)
}

func TestTodoHandler_GetByID(t *testing.T) {
	svc := mocks.NewMockService(t)
	handler := NewTodoHandler(svc)

	expectedTodo := &todo.Todo{
		ID:    7,
		Title: "Task",
	}

	svc.On("GetById", mock.Anything, int64(7)).Return(expectedTodo, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/todos/7", nil)
	req = withURLParam(req, "id", "7")
	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var got todo.Todo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, *expectedTodo, got)

	svc.AssertExpectations(t)
}

func TestTodoHandler_GetByID_InvalidID(t *testing.T) {
	handler := NewTodoHandler(mocks.NewMockService(t))

	req := httptest.NewRequest(http.MethodGet, "/todos/abc", nil)
	req = withURLParam(req, "id", "abc")
	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid id")
}

func TestTodoHandler_GetAll(t *testing.T) {
	svc := mocks.NewMockService(t)
	handler := NewTodoHandler(svc)

	priority := todo.High
	expectedTodos := []todo.Todo{
		{ID: 1, Title: "Task 1", Priority: todo.High},
	}

	svc.On("GetAll", mock.Anything, &priority).Return(expectedTodos, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/todos?priority=high", nil)
	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var got []todo.Todo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, expectedTodos, got)

	svc.AssertExpectations(t)
}

func TestTodoHandler_GetAll_WithoutPriority(t *testing.T) {
	svc := mocks.NewMockService(t)
	handler := NewTodoHandler(svc)

	expectedTodos := []todo.Todo{{ID: 1, Title: "Task 1"}}

	svc.On("GetAll", mock.Anything, (*todo.Priority)(nil)).Return(expectedTodos, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	svc.AssertExpectations(t)
}

func TestTodoHandler_Update(t *testing.T) {
	svc := mocks.NewMockService(t)
	handler := NewTodoHandler(svc)

	priority := todo.Medium
	completed := true
	expectedReq := &todo.UpdateTodoRequest{
		Title:     "Updated Task",
		Priority:  &priority,
		Completed: &completed,
	}
	expectedTodo := &todo.Todo{
		ID:        1,
		Title:     "Updated Task",
		Priority:  todo.Medium,
		Completed: true,
	}

	svc.On("Update", mock.Anything, int64(1), expectedReq).Return(expectedTodo, nil).Once()

	req := httptest.NewRequest(http.MethodPatch, "/todos/1", bytes.NewBufferString(`{"title":"Updated Task","priority":"medium","completed":true}`))
	req = withURLParam(req, "id", "1")
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var got todo.Todo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, *expectedTodo, got)

	svc.AssertExpectations(t)
}

func TestTodoHandler_Update_InvalidID(t *testing.T) {
	handler := NewTodoHandler(mocks.NewMockService(t))

	req := httptest.NewRequest(http.MethodPatch, "/todos/abc", bytes.NewBufferString(`{"title":"Updated Task"}`))
	req = withURLParam(req, "id", "abc")
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid id")
}

func TestTodoHandler_Update_InvalidBody(t *testing.T) {
	handler := NewTodoHandler(mocks.NewMockService(t))

	req := httptest.NewRequest(http.MethodPatch, "/todos/1", bytes.NewBufferString("{"))
	req = withURLParam(req, "id", "1")
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid body")
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code int
		body string
	}{
		{
			name: "not found",
			err:  todo.ErrTodoNotFound,
			code: http.StatusNotFound,
			body: `{"error":"todo not found"}`,
		},
		{
			name: "invalid priority",
			err:  todo.ErrInvalidPriority,
			code: http.StatusBadRequest,
			body: `{"error":"invalid priority"}`,
		},
		{
			name: "internal error",
			err:  errors.New("boom"),
			code: http.StatusInternalServerError,
			body: `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			handleError(rec, tt.err)

			require.Equal(t, tt.code, rec.Code)
			require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			require.JSONEq(t, tt.body, rec.Body.String())
		})
	}
}
