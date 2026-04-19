package postgres

import (
	"context"
	"fmt"
	"task-tracking/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NamespaceRepo struct {
	pool *pgxpool.Pool
}

func NewNamespaceRepo(pool *pgxpool.Pool) *NamespaceRepo {
	return &NamespaceRepo{pool: pool}
}

func (r *NamespaceRepo) Create(ctx context.Context, namespace domain.Namespace) (*domain.Namespace, error) {
	query := `
		INSERT INTO namespaces (name)
		VALUES ($1)
		RETURNING id, name, created_at
	`
	err := r.pool.QueryRow(ctx, query,
		namespace.Name,
	).Scan(&namespace.ID, &namespace.Name, &namespace.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create namespace: %w", err)
	}

	return &namespace, nil
}

func (r *NamespaceRepo) Update(ctx context.Context, namespace domain.Namespace) (*domain.Namespace, error) {
	query := `
		UPDATE namespaces
		SET name = $1
		WHERE id = $2
		RETURNING id, name, updated_at`

	err := r.pool.QueryRow(ctx, query,
		namespace.Name, namespace.ID,
	).Scan(&namespace.ID, &namespace.Name, &namespace.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("update namespace: %w", err)
	}

	return &namespace, nil
}

func (r *NamespaceRepo) GetByNamespaceId(ctx context.Context, namespaceId int64) (*domain.Namespace, error) {
	var namespace domain.Namespace
	query := `
		SELECT id, name, created_at, updated_at
		FROM namespaces
		WHERE id = $1`

	err := r.pool.QueryRow(ctx, query,
		namespaceId).Scan(&namespace.ID, &namespace.Name, &namespace.CreatedAt, &namespace.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("get namespace by id: %w", err)
	}

	return &namespace, nil
}

func (r *NamespaceRepo) GetByNamespaceName(ctx context.Context, namespaceName string) (*domain.Namespace, error) {
	var namespace domain.Namespace
	query := `
		SELECT id, name, created_at, updated_at
		FROM namespaces
		WHERE name = $1`

	err := r.pool.QueryRow(ctx, query, namespaceName).Scan(&namespace.ID, &namespace.Name, &namespace.CreatedAt, &namespace.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("get namespace by name: %w", err)
	}

	return &namespace, nil
}
