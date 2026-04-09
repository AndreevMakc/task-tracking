package repository

import "task-tracking/internal/domain"

type NamespaceRepository interface {
	Save(namespace domain.Namespace) error
	FindByName(name string) *domain.Namespace
}
