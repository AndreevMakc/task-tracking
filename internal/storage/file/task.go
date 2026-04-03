package file

import (
	"encoding/json"
	"fmt"
	"os"
	"task-tracking/internal/domain"
)

const DefaultNamespaceName = "default"

type NamespaceRepository struct {
	filePath string
	storage  StorageRecord
}

func toDomain(namespaceRecord NamespaceRecord) domain.Namespace {
	var tasks []domain.Task
	for index := range namespaceRecord.Tasks {
		tasks = append(tasks,
			domain.Task{
				ID:     namespaceRecord.Tasks[index].ID,
				SeqId:  namespaceRecord.Tasks[index].SeqId,
				Title:  namespaceRecord.Tasks[index].Title,
				IsDone: namespaceRecord.Tasks[index].IsDone})
	}
	return domain.Namespace{Name: namespaceRecord.Namespace, Counter: namespaceRecord.Counter, Tasks: tasks}
}

func toRecord(namespace domain.Namespace) NamespaceRecord {
	var tasks []TaskRecord
	for index := range namespace.Tasks {
		tasks = append(tasks,
			TaskRecord{
				ID:     namespace.Tasks[index].ID,
				SeqId:  namespace.Tasks[index].SeqId,
				Title:  namespace.Tasks[index].Title,
				IsDone: namespace.Tasks[index].IsDone,
			})
	}
	return NamespaceRecord{Namespace: namespace.Name, Counter: namespace.Counter, Tasks: tasks}
}

func (f *NamespaceRepository) FindByName(namespaceName string) (*domain.Namespace, error) {

	namespaceDomain := domain.Namespace{}
	for index := range f.storage.Items {
		if f.storage.Items[index].Namespace == namespaceName {
			namespaceDomain = toDomain(f.storage.Items[index])
			return &namespaceDomain, nil
		}
	}
	return nil, nil
}

func (f *NamespaceRepository) Save(namespace domain.Namespace) error {
	found := false
	for index := range f.storage.Items {
		if f.storage.Items[index].Namespace == namespace.Name {
			f.storage.Items[index] = toRecord(namespace)
			found = true
			break
		}
	}
	if !found {
		f.storage.Items = append(f.storage.Items, toRecord(namespace))
	}
	data, err := json.MarshalIndent(f.storage, "", "")
	if err != nil {
		return fmt.Errorf("serialization error: %w", err)
	}
	if err := os.WriteFile(f.filePath, data, 0644); err != nil {
		return fmt.Errorf("file write error: %w", err)
	}
	return nil
}

func NewNamespaceRepository(filePath string) (*NamespaceRepository, error) {
	namespaceRepository := NamespaceRepository{filePath: filePath}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &namespaceRepository, nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("file read error: %w", err)
	}
	if err := json.Unmarshal(data, &namespaceRepository.storage); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %w", err)
	}

	return &namespaceRepository, nil
}
