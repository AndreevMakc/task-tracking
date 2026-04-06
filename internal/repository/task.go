package repository

import "task-tracking/internal/domain"

type TaskRepository interface {
	Save(task domain.Task, namespaceName string) error
	FindBySeqId(seqId int, namespaceName string) (*domain.Task, error)
	FindAllByNamespaceName(namespaceName string) ([]domain.Task, error)
}

type NamespaceRepository interface {
	Save(namespace domain.Namespace) error
	FindByName(name string) (*domain.Namespace, error)
}
