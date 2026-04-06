package file

type StorageRecord struct {
	Items []NamespaceRecord `json:"items"`
}

type NamespaceRecord struct {
	Namespace string       `json:"namespace"`
	Counter   int          `json:"counter"`
	Tasks     []TaskRecord `json:"tasks"`
}

type TaskRecord struct {
	ID     string `json:"id"`
	SeqId  int    `json:"seq_id"`
	Title  string `json:"title"`
	IsDone bool   `json:"is_done"`
}
