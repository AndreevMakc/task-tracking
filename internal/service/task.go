package service

import (
	"fmt"
	"strconv"
	"strings"
	"task-tracking/internal/domain"
	"task-tracking/internal/repository"

	"github.com/google/uuid"
)

type NamespaceService struct {
	namespaceRepo repository.NamespaceRepository
}

func NewNamespaceService(repo repository.NamespaceRepository) *NamespaceService {
	return &NamespaceService{namespaceRepo: repo}
}

type TaskService struct {
	namespaceRepository repository.NamespaceRepository
}

func NewTaskService(repo repository.NamespaceRepository) *TaskService {
	return &TaskService{namespaceRepository: repo}
}

func (s *TaskService) CreateTask(text string, namespaceName string) (domain.Task, error) {
	if text == "" {
		return domain.Task{}, fmt.Errorf("task title is required")
	}
	namespace := s.namespaceRepository.FindByName(namespaceName)
	if namespace == nil {
		namespace = &domain.Namespace{Name: namespaceName}
	}
	task := domain.Task{ID: uuid.New().String(), SeqId: len(namespace.Tasks) + 1, Title: text, IsDone: false}
	namespace.Tasks = append(namespace.Tasks, task)
	if err := s.namespaceRepository.Save(*namespace); err != nil {
		return domain.Task{}, fmt.Errorf("task save error, %w", err)
	}
	return task, nil
}

func (s *TaskService) TaskDone(humanId string, namespaceName string) error {
	namespace := s.namespaceRepository.FindByName(namespaceName)
	lastSepIndex := strings.LastIndex(humanId, "-")
	seqId, err := strconv.Atoi(humanId[lastSepIndex+1:])
	if err != nil {
		return fmt.Errorf("invalid task id: %s", humanId)
	}
	found := false
	for i := range namespace.Tasks {
		if namespace.Tasks[i].SeqId == seqId {
			namespace.Tasks[i].IsDone = true
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("task not found: %s", humanId)
	}
	return s.namespaceRepository.Save(*namespace)
}

func (s *TaskService) ShowList(namespaceName string) []string {
	namespace := s.namespaceRepository.FindByName(namespaceName)
	var tasks []string
	if namespace != nil {
		for i := range namespace.Tasks {
			tasks = append(tasks, fmt.Sprintf("[%s-%d] %s [%t]",
				namespace.Name,
				namespace.Tasks[i].SeqId,
				namespace.Tasks[i].Title,
				namespace.Tasks[i].IsDone))
		}
	}
	if len(tasks) == 0 {
		tasks = append(tasks, "No tasks yet")
	}
	return tasks
}

func (s *TaskService) GetTasks(namespaceName string) []domain.Task {
	namespace := s.namespaceRepository.FindByName(namespaceName)
	if namespace == nil {
		return []domain.Task{}
	}
	return namespace.Tasks
}
