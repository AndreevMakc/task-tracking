package filestorage

import (
	"encoding/json"
	"fmt"
	"os"
	"task-tracking/internal/domain"
)

const defaultFilePath = "tasks.json"

type FileStorage struct {
	storage  StorageDTO
	filePath string
}

func (f *FileStorage) loadFile() error {
	if f.filePath == "" {
		f.filePath = defaultFilePath
	}
	if _, err := os.Stat(f.filePath); os.IsNotExist(err) {
		return nil
	}
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return fmt.Errorf("file storage read error: %w", err)
	}
	if err := json.Unmarshal(data, &f.storage); err != nil {
		return fmt.Errorf("JSON parsing error: %w", err)
	}
	return nil
}

func (f *FileStorage) saveFile() error {
	data, err := json.MarshalIndent(f.storage, "", "")
	if err != nil {
		return fmt.Errorf("serialization error: %w", err)
	}
	if err := os.WriteFile(f.filePath, data, 0644); err != nil {
		return fmt.Errorf("file storage write error: %w", err)
	}
	return nil
}

func mapDTOTONamespace(namespaceRecord NamespaceDTO) domain.Namespace {
	var tasks []domain.Task
	for index := range namespaceRecord.Tasks {
		tasks = append(tasks,
			domain.Task{
				ID:     namespaceRecord.Tasks[index].ID,
				SeqId:  namespaceRecord.Tasks[index].SeqId,
				Title:  namespaceRecord.Tasks[index].Title,
				IsDone: namespaceRecord.Tasks[index].IsDone})
	}
	return domain.Namespace{Name: namespaceRecord.Namespace, Tasks: tasks}
}

func mapNamespaceTODTO(namespace domain.Namespace) NamespaceDTO {
	var tasks []TaskDTO
	for index := range namespace.Tasks {
		tasks = append(tasks,
			TaskDTO{
				ID:     namespace.Tasks[index].ID,
				SeqId:  namespace.Tasks[index].SeqId,
				Title:  namespace.Tasks[index].Title,
				IsDone: namespace.Tasks[index].IsDone,
			})
	}
	return NamespaceDTO{Namespace: namespace.Name, Tasks: tasks}
}

type FileNamespaceRepository struct {
	fileStorage *FileStorage
}

func (f *FileNamespaceRepository) Save(namespace domain.Namespace) error {
	found := false
	for i := range f.fileStorage.storage.Items {
		if f.fileStorage.storage.Items[i].Namespace == namespace.Name {
			f.fileStorage.storage.Items[i] = mapNamespaceTODTO(namespace)
			found = true
			break
		}
	}
	if !found {
		f.fileStorage.storage.Items = append(f.fileStorage.storage.Items, mapNamespaceTODTO(namespace))
	}
	if err := f.fileStorage.saveFile(); err != nil {
		return fmt.Errorf("file save error: %w", err)
	}
	return nil
}

func (f *FileNamespaceRepository) FindByName(namespace string) *domain.Namespace {
	if namespace == "" {
		namespace = domain.DefaultNamespace
	}
	for i := range f.fileStorage.storage.Items {
		if f.fileStorage.storage.Items[i].Namespace == namespace {
			namespaceDomain := mapDTOTONamespace(f.fileStorage.storage.Items[i])
			return &namespaceDomain
		}
	}
	return nil
}

func NewFileNamespaceRepository(filePath string) (*FileNamespaceRepository, error) {
	fileStorage := FileStorage{filePath: filePath}
	if err := fileStorage.loadFile(); err != nil {
		return nil, fmt.Errorf("file load error: %w", err)
	}
	return &FileNamespaceRepository{fileStorage: &fileStorage}, nil
}
