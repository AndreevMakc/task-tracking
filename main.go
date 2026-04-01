package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

const defaultTenant = "my-tenant"

type App struct {
	Storage Storage
	Current *TenantTask
}

type Storage struct {
	Items []TenantTask `json:"items"`
}

type TenantTask struct {
	Tenant  string `json:"tenant"`
	Counter int    `json:"counter"`
	Tasks   []Task `json:"tasks"`
}

type Task struct {
	ID     string `json:"id"`
	SeqId  int    `json:"seq_id"`
	Title  string `json:"title"`
	IsDone bool   `json:"is_done"`
}

type TaskView struct {
	HumanId string `json:"id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
}

func newApp() (*App, error) {
	storage, err := loadFile()
	if err != nil {
		return nil, fmt.Errorf("app init failed: %w", err)
	}
	var currentTenant string
	if len(os.Args) > 1 {
		currentTenant = os.Args[1]
	}
	currentTask, err := storage.findOrCreateTaskTenant(currentTenant)
	if err != nil {
		return nil, fmt.Errorf("app init failed: %w", err)
	}
	return &App{Storage: storage, Current: currentTask}, nil
}

func loadFile() (Storage, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return Storage{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return Storage{}, fmt.Errorf("file read error: %w", err)
	}

	var storage Storage
	if err := json.Unmarshal(data, &storage); err != nil {
		return Storage{}, fmt.Errorf("JSON parsing error: %w", err)
	}

	return storage, nil
}

func (storage *Storage) findOrCreateTaskTenant(tenant string) (*TenantTask, error) {
	if tenant == "" {
		tenant = defaultTenant
	}
	for index := range storage.Items {
		if storage.Items[index].Tenant == tenant {
			return &storage.Items[index], nil
		}
	}
	storage.Items = append(storage.Items, TenantTask{Tenant: tenant})
	return &storage.Items[len(storage.Items)-1], nil
}

func getAllTasks(tenantTask *TenantTask) ([]Task, error) {
	return tenantTask.Tasks, nil
}

func humanId(tenant string, seqId int) string {
	return fmt.Sprintf("%s-%d", tenant, seqId)
}

func genUUID() string {
	return uuid.New().String()
}

func createTask(text string, tenantTask *TenantTask) (Task, error) {
	if text == "" {
		return Task{}, fmt.Errorf("task title is required")
	}
	tenantTask.Counter = tenantTask.Counter + 1
	return Task{ID: genUUID(), SeqId: tenantTask.Counter, Title: text, IsDone: false}, nil
}

func handleInputCmd(text string) (cmd string, taskTitle string) {
	sliceText := strings.SplitN(text, " ", 2)
	cmd = sliceText[0]
	if len(sliceText) > 1 {
		taskTitle = sliceText[1]
	}
	return cmd, taskTitle
}

const filePath = "tasks.json"

func saveFile(storage Storage) error {
	data, err := json.MarshalIndent(storage, "", " ")
	if err != nil {
		return fmt.Errorf("serialization error: %w", err)
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("file write error: %w", err)
	}
	return nil
}

func taskDone(taskId string, tasks []Task) ([]Task, error) {
	for i := range tasks {
		if tasks[i].ID == taskId {
			tasks[i].IsDone = true
			return tasks, nil
		}
	}
	return []Task{}, fmt.Errorf("task not found: %s", taskId)
}

func main() {
	app, err := newApp()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("---Task tracking is running---")
	fmt.Println("Type 'add', 'list', 'done' or 'exit'")
	for {
		fmt.Print(">")
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		cmd, params := handleInputCmd(text)
		switch cmd {
		case "exit":
			if err := saveFile(app.Storage); err != nil {
				fmt.Fprintf(os.Stderr, "File save error: %s\n", err)
			}
			fmt.Println("Goodbye")
			return
		case "add":
			task, err := createTask(params, app.Current)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Task add error: %s\n", err)
				continue
			}
			app.Current.Tasks = append(app.Current.Tasks, task)
		case "list":
			if len(app.Current.Tasks) == 0 {
				fmt.Println("No tasks yet")
				continue
			}
			for _, v := range app.Current.Tasks {
				fmt.Printf("[%s] %s - %t\n", v.ID, v.Title, v.IsDone)
			}
		}
	}
}
