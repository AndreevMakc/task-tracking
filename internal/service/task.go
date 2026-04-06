package service

import (
	"fmt"
	"strconv"
	"strings"
	"task-tracking/internal/domain"

	"github.com/google/uuid"
)

func CreateTask(text string, namespace *domain.Namespace) (domain.Task, error) {
	if text == "" {
		return domain.Task{}, fmt.Errorf("task title is required")
	}
	namespace.Counter = namespace.Counter + 1
	return domain.Task{ID: uuid.New().String(), SeqId: namespace.Counter, Title: text, IsDone: false}, nil
}

func TaskDone(humanId string, namespace *domain.Namespace) error {
	lastSepIndex := strings.LastIndex(humanId, "-")
	seqId, err := strconv.Atoi(humanId[lastSepIndex+1:])
	if err != nil {
		return fmt.Errorf("invalid task id: %s", humanId)
	}
	for i := range namespace.Tasks {
		if namespace.Tasks[i].SeqId == seqId {
			namespace.Tasks[i].IsDone = true
			return nil
		}
	}
	return fmt.Errorf("task not found: %s", humanId)
}
