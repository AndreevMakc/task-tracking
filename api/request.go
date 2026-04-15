package api

type CreateTaskRequest struct {
	Title string `json:"title"`
}

type PatchTaskRequest struct {
	IsDone bool `json:"is_done"`
}
