package domain

type TaskStatus string

const (
	TaskStatusNew        TaskStatus = "New"
	TaskStatusInProgress TaskStatus = "InProgress"
	TaskStatusDone       TaskStatus = "Done"
	TaskStatusTrashed    TaskStatus = "Trashed"
)
