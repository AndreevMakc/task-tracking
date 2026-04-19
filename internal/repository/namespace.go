package repository

import (
	"context"
	"task-tracking/internal/domain"
)

type NamespaceRepository interface {
	Create(ctx context.Context, namespace domain.Namespace) (*domain.Namespace, error)
	Update(ctx context.Context, namespace domain.Namespace) (*domain.Namespace, error)
	GetByNamespaceId(ctx context.Context, namespaceId int64) (*domain.Namespace, error)
	GetByNamespaceName(ctx context.Context, namespaceName string) (*domain.Namespace, error)
}
