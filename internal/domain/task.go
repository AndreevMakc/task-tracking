package domain

type Task struct {
	ID     string
	SeqId  int
	Title  string
	IsDone bool
}

type Namespace struct {
	Name    string
	Counter int
	Tasks   []Task
}
