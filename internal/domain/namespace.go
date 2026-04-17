package domain

const DefaultNamespace = "default"

type Namespace struct {
	BaseModelID
	Name string
}
