package domain

const DefaultNamespace = "default"

type Namespace struct {
	Name  string
	Tasks []Task
}
