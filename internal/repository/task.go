package repository

import (
	"context"

	"task-tracking/internal/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, task domain.Task) (*domain.Task, error)
	Update(ctx context.Context, task domain.Task) (*domain.Task, error)
	GetByNamespaceId(ctx context.Context, namespaceId int64) ([]domain.Task, error)
	GetByCode(ctx context.Context, code string) (*domain.Task, error)
}
