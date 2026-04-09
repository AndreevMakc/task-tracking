package filestorage

type StorageDTO struct {
	Items []NamespaceDTO `json:"items"`
}

type NamespaceDTO struct {
	Namespace string    `json:"namespace"`
	Tasks     []TaskDTO `json:"tasks"`
}

type TaskDTO struct {
	ID     string `json:"id"`
	SeqId  int    `json:"seq_id"`
	Title  string `json:"title"`
	IsDone bool   `json:"is_done"`
}
