package service

import (
	"context"
	"fmt"

	"task-tracking/internal/domain"
	"task-tracking/internal/repository"
)

type TaskService struct {
	namespaceRepo repository.NamespaceRepository
	taskRepo      repository.TaskRepository
}

func NewTaskService(
	namespaceRepo repository.NamespaceRepository,
	taskRepo repository.TaskRepository,
) *TaskService {
	return &TaskService{namespaceRepo: namespaceRepo, taskRepo: taskRepo}
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	title string,
	namespaceName string,
) (*domain.Task, error) {
	namespaceId, err := s.resolveNamespaceId(ctx, namespaceName)
	if err != nil {
		return nil, fmt.Errorf("resolve namespace: %w", err)
	}
	task, err := s.taskRepo.Create(ctx, domain.Task{Title: title, NamespaceID: namespaceId})
	if err != nil {
		return nil, fmt.Errorf("task create: %w", err)
	}
	return task, nil
}

func (s *TaskService) TaskDone(ctx context.Context, code string) (*domain.Task, error) {
	task, err := s.changeTaskStatus(ctx, code, domain.TaskStatusDone)
	if err != nil {
		return nil, fmt.Errorf("task done: %w", err)
	}
	return task, nil
}

func (s *TaskService) TaskTrash(ctx context.Context, code string) (*domain.Task, error) {
	task, err := s.changeTaskStatus(ctx, code, domain.TaskStatusTrashed)
	if err != nil {
		return nil, fmt.Errorf("task trash: %w", err)
	}
	return task, nil
}

func (s *TaskService) GetTasksByNamespace(
	ctx context.Context,
	namespaceName string,
) ([]domain.Task, error) {
	namespace, err := s.namespaceRepo.GetByNamespaceName(ctx, namespaceName)
	if err != nil {
		return nil, fmt.Errorf("get namespace by name: %w", err)
	}
	if namespace == nil {
		return nil, fmt.Errorf("namespace %s not found", namespaceName)
	}
	tasks, err := s.taskRepo.GetByNamespaceId(ctx, namespace.ID)
	if err != nil {
		return nil, fmt.Errorf("get tasks by namespace id: %w", err)
	}
	return tasks, nil
}

func (s *TaskService) resolveNamespaceId(ctx context.Context, namespaceName string) (int64, error) {
	namespace, err := s.namespaceRepo.GetByNamespaceName(ctx, namespaceName)
	if err != nil {
		return 0, fmt.Errorf("get namespace by name: %w", err)
	}
	if namespace == nil {
		namespace, err = s.namespaceRepo.GetByNamespaceName(ctx, domain.DefaultNamespace)
		if err != nil {
			return 0, fmt.Errorf("get default namespace by name: %w", err)
		}
		if namespace == nil {
			return 0, fmt.Errorf("default namespace not found")
		}
	}

	return namespace.ID, nil
}

func (s *TaskService) changeTaskStatus(
	ctx context.Context,
	code string,
	status domain.TaskStatus,
) (*domain.Task, error) {
	task, err := s.taskRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get task by code: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task %s not found", code)
	}
	task.Status = status
	task, err = s.taskRepo.Update(ctx, *task)
	if err != nil {
		return nil, fmt.Errorf("task update: %w", err)
	}
	return task, nil
}
