package main

import (
	"fmt"
	"net/http"
	"os"
	"task-tracking/api"
	"task-tracking/internal/service"
	"task-tracking/internal/storage/filestorage"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.elastic.co/apm/module/apmchiv5/v2"
	"go.elastic.co/apm/v2"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(apmchiv5.Middleware())
	fmt.Println("APM URL:", os.Getenv("ELASTIC_APM_SERVER_URL"))
	fmt.Println("APM Token set:", os.Getenv("ELASTIC_APM_SECRET_TOKEN") != "")
	fmt.Println("APM Service:", os.Getenv("ELASTIC_APM_SERVICE_NAME"))
	fileNamespaceRepo, err := filestorage.NewFileNamespaceRepository("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	taskService := service.NewTaskService(fileNamespaceRepo)
	namespaceHandler := api.NewNamespaceHandler(taskService)
	r.Get("/namespaces/{namespace}/tasks", namespaceHandler.GetTasks)
	r.Post("/namespaces/{namespace}/tasks", namespaceHandler.CreateTask)
	r.Patch("/namespaces/{namespace}/tasks/{taskId}", namespaceHandler.TaskDone)
	defer apm.DefaultTracer().Flush(nil)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
