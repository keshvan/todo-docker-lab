package todo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{db}
}

func (r *TodoRepository) Create(ctx context.Context, todo *Todo) error {
	query := `
		INSERT INTO todos (title, description, priority, completed, deadline)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query, todo.Title, todo.Description, todo.Priority, todo.Completed, todo.Deadline).Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt)
}

func (r *TodoRepository) GetById(ctx context.Context, id int64) (*Todo, error) {
	query := `
		SELECT * FROM todos WHERE id = $1
	`

	var todo Todo
	err := r.db.QueryRow(ctx, query, id).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.Priority,
		&todo.Completed, &todo.Deadline, &todo.CreatedAt, &todo.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepository) GetAll(ctx context.Context, priority *Priority) ([]Todo, error) {

	query := `
		SELECT * FROM todos
	`

	args := []any{}
	if priority != nil {
		query += " WHERE priority = $1"
		args = append(args, *priority)
	}

	query += `
		ORDER BY 
			CASE priority
				WHEN 'high' THEN 1
				WHEN 'medium' THEN 2
				WHEN 'low' THEN 3
			END,
		created_at DESC
	`

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos = make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Description, &todo.Priority,
			&todo.Completed, &todo.Deadline, &todo.CreatedAt, &todo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, rows.Err()
}

func (r *TodoRepository) Update(ctx context.Context, id int64, todo *Todo) error {
	query := `
		UPDATE todos
		SET title = $1, description = $2, priority = $3, completed = $4, deadline = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(ctx, query, todo.Title, todo.Description, todo.Priority, todo.Completed, todo.Deadline, id)
	return err
}
