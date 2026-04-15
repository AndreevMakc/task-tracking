package api

import (
	"encoding/json"
	"net/http"
	"task-tracking/internal/service"

	"github.com/go-chi/chi/v5"
)

type NamespaceHandler struct {
	taskService *service.TaskService
}

func NewNamespaceHandler(taskService *service.TaskService) *NamespaceHandler {
	return &NamespaceHandler{taskService: taskService}
}

func (h *NamespaceHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	if namespace == "" {
		http.Error(w, "missing namespace", http.StatusBadRequest)
		return
	}
	tasks := h.taskService.GetTasks(namespace)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *NamespaceHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	namespace := chi.URLParam(r, "namespace")
	if namespace == "" {
		http.Error(w, "missing namespace", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	task, err := h.taskService.CreateTask(req.Title, namespace)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		return
	}
}

func (h *NamespaceHandler) TaskDone(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	humanId := chi.URLParam(r, "taskId")
	if namespace == "" {
		http.Error(w, "missing namespace", http.StatusBadRequest)
		return
	}
	if humanId == "" {
		http.Error(w, "missing humanId", http.StatusBadRequest)
		return
	}
	if err := h.taskService.TaskDone(humanId, namespace); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
