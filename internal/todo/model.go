package todo

import "time"

type Priority string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

type Todo struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    Priority   `json:"priority"`
	Completed   bool       `json:"completed"`
	Deadline    *time.Time `json:"deadline"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTodoRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Priority    Priority   `json:"priority"`
	Deadline    *time.Time `json:"deadline,omitempty"`
}

type UpdateTodoRequest struct {
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Priority    *Priority `json:"priority,omitempty"`
	Completed   *bool     `json:"completed,omitempty"`
}

func (r *CreateTodoRequest) SetDefaults() {
	if r.Priority == "" {
		r.Priority = Medium
	}
}
