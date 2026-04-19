package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"task-tracking/internal/domain"
)

type TaskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{pool: pool}
}

func (r *TaskRepo) Create(ctx context.Context, task domain.Task) (*domain.Task, error) {
	query := `
		INSERT INTO tasks (namespace_id, title)
		VALUES ($1, $2)
		RETURNING id, code, title, status, namespace_id, created_at`

	err := r.pool.QueryRow(ctx, query,
		task.NamespaceID, task.Title,
	).Scan(&task.ID, &task.Code, &task.Title, &task.Status, &task.NamespaceID, &task.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return &task, nil
}

func (r *TaskRepo) Update(ctx context.Context, task domain.Task) (*domain.Task, error) {
	query := `
		UPDATE tasks 
		SET title = $1, status = $2
		WHERE code = $3
		RETURNING id, code, title, status, updated_at`

	err := r.pool.QueryRow(ctx, query,
		task.Title, task.Status,
		task.Code,
	).Scan(&task.ID, &task.Code, &task.Title, &task.Status, &task.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	return &task, nil
}

func (r *TaskRepo) GetByNamespaceId(ctx context.Context, namespaceId int64) ([]domain.Task, error) {
	query := `
		SELECT id, namespace_id, title, code, status, created_at, updated_at, deleted_at
		FROM tasks
		WHERE namespace_id = $1`

	rows, err := r.pool.Query(ctx, query, namespaceId)
	if err != nil {
		return nil, fmt.Errorf("list namespace tasks: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task

		if err := rows.Scan(
			&task.ID, &task.NamespaceID, &task.Title, &task.Code,
			&task.Status, &task.CreatedAt, &task.UpdatedAt, &task.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepo) GetByCode(ctx context.Context, code string) (*domain.Task, error) {
	var task domain.Task
	query := `
		SELECT id, namespace_id, title, code, status, created_at, updated_at, deleted_at
		FROM tasks
		WHERE code = $1`

	err := r.pool.QueryRow(ctx, query, code).Scan(
		&task.ID, &task.NamespaceID,
		&task.Title, &task.Code,
		&task.Status, &task.CreatedAt,
		&task.UpdatedAt, &task.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get task by code: %w", err)
	}

	return &task, nil
}
