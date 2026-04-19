package domain

type Task struct {
	BaseModelUUID
	Code        string
	Title       string
	Status      TaskStatus
	NamespaceID int64
}
