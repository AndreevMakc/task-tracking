package domain

type Task struct {
	BaseModelUUID
	Title  string
	Status TaskStatus
}
