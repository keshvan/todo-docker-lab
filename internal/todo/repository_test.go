package todo_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	todo "github.com/keshvan/todo-docker-lab/internal/todo"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := todo.NewTodoRepository(mock)

	desc := "desc"
	deadline := time.Now()

	todo := &todo.Todo{
		Title:       "test",
		Description: &desc,
		Priority:    todo.High,
		Completed:   false,
		Deadline:    &deadline,
	}

	rows := pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, time.Now(), time.Now())

	mock.ExpectQuery("INSERT INTO todos").
		WithArgs(todo.Title, todo.Description, todo.Priority, todo.Completed, todo.Deadline).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), todo)

	require.NoError(t, err)
	require.Equal(t, int64(1), todo.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetById_Found(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := todo.NewTodoRepository(mock)

	desc := "desc"
	deadline := time.Now()
	created := time.Now()
	updated := time.Now()

	rows := pgxmock.NewRows([]string{
		"id", "title", "description", "priority",
		"completed", "deadline", "created_at", "updated_at",
	}).AddRow(1, "test", &desc, todo.High, false, &deadline, created, updated)

	mock.ExpectQuery("SELECT \\* FROM todos WHERE id = \\$1").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	todo, err := repo.GetById(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, todo)
	require.Equal(t, int64(1), todo.ID)
	require.Equal(t, &desc, todo.Description)
	require.Equal(t, &deadline, todo.Deadline)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetById_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := todo.NewTodoRepository(mock)

	mock.ExpectQuery("SELECT \\* FROM todos WHERE id = \\$1").
		WithArgs(int64(1)).
		WillReturnError(pgx.ErrNoRows)

	todo, err := repo.GetById(context.Background(), 1)

	require.NoError(t, err)
	require.Nil(t, todo)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := todo.NewTodoRepository(mock)

	desc1 := "d1"
	desc2 := "d2"
	deadline1 := time.Now()
	deadline2 := time.Now()
	created1 := time.Now()
	created2 := time.Now()
	updated1 := time.Now()
	updated2 := time.Now()

	rows := pgxmock.NewRows([]string{
		"id", "title", "description", "priority",
		"completed", "deadline", "created_at", "updated_at",
	}).
		AddRow(1, "t1", &desc1, todo.High, false, &deadline1, created1, updated1).
		AddRow(2, "t2", &desc2, todo.Low, true, &deadline2, created2, updated2)

	mock.ExpectQuery("SELECT \\* FROM todos").
		WillReturnRows(rows)

	todos, err := repo.GetAll(context.Background(), nil)

	require.NoError(t, err)
	require.Len(t, todos, 2)
	require.Equal(t, &desc1, todos[0].Description)
	require.Equal(t, &desc2, todos[1].Description)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_WithPriority(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := todo.NewTodoRepository(mock)

	desc := "task description"
	deadline := time.Now()
	created := time.Now()
	updated := time.Now()
	priority := todo.High

	rows := pgxmock.NewRows([]string{
		"id", "title", "description", "priority",
		"completed", "deadline", "created_at", "updated_at",
	}).AddRow(1, "Task1", &desc, priority, false, &deadline, created, updated)

	mock.ExpectQuery("SELECT \\* FROM todos WHERE priority = \\$1").
		WithArgs(priority).
		WillReturnRows(rows)

	todos, err := repo.GetAll(context.Background(), &priority)

	require.NoError(t, err)
	require.Len(t, todos, 1)
	require.Equal(t, &desc, todos[0].Description)
	require.Equal(t, priority, todos[0].Priority)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := todo.NewTodoRepository(mock)

	desc := "Updated description"
	deadline := time.Now()
	todo := &todo.Todo{
		Title:       "Updated Task",
		Description: &desc,
		Priority:    todo.Medium,
		Completed:   true,
		Deadline:    &deadline,
	}

	mock.ExpectExec("UPDATE todos").
		WithArgs(todo.Title, todo.Description, todo.Priority, todo.Completed, todo.Deadline, int64(1)).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), 1, todo)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
