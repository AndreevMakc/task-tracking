package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

const tenant = "my-tenant"

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

func (app App) init() (Storage, error) {
	storage, err := loadFile()
	if err != nil {
		panic(fmt.Sprint("app init failed:", err))
	}
	return storage, nil
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
	for index := range storage.Items {
		if storage.Items[index].Tenant == tenant {
			return &storage.Items[index], nil
		}
	}
	storage.Items = append(storage.Items, TenantTask{Tenant: tenant})
	return &storage.Items[len(storage.Items)-1], nil
}

func humanId(tenant string, seqId int) string {
	return fmt.Sprintf("%s-%d", tenant, seqId)
}

func genUUID() string {
	return uuid.New().String()
}

func createTask(text string) (Task, error) {
	if text == "" {
		return Task{}, fmt.Errorf("task title is required")
	}
	return Task{ID: genUUID(), Title: text, IsDone: false}, nil
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

func saveFile(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
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
	var app App
	storage, err := app.init()
	if err != nil {
		panic(err)
	}
	tasks, _ := storage.findOrCreateTaskTenant(tenant)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("---Task tracking is running---")
	fmt.Println("Type 'add', 'list', 'done' or 'exit'")
	for i := range tasks.Tasks {
		fmt.Println(humanId(tasks.Tenant, tasks.Tasks[i].SeqId))
	}
	// for {
	// 	fmt.Print("> ")
	// 	if !scanner.Scan() {
	// 		break
	// 	}
	// 	text := scanner.Text()
	// 	cmd, params := handleInputCmd(text)
	// 	switch cmd {
	// 	case "exit":
	// 		if err := saveFile(tasks); err != nil {
	// 			fmt.Fprintf(os.Stderr, "File save error: %s\n", err)
	// 		}
	// 		fmt.Println("Goodbye")
	// 		return
	// 	case "add":
	// 		task, err := createTask(params)
	// 		if err != nil {
	// 			fmt.Fprintf(os.Stderr, "Task add error: %s\n", err)
	// 			continue
	// 		}
	// 		tasks = append(tasks, task)
	// 		fmt.Printf("Task '%s' has been added\n", task.Title)
	// 	case "list":
	// 		if len(tasks) == 0 {
	// 			fmt.Println("No tasks yet")
	// 			continue
	// 		}
	// 		for _, v := range tasks {
	// 			fmt.Printf("[%s] %s - %t\n", v.ID, v.Title, v.IsDone)
	// 		}
	// 	case "done":
	// 		updatedTasks, err := taskDone(params, tasks)
	// 		if err != nil {
	// 			fmt.Fprintf(os.Stderr, "Task done error: %s\n", err)
	// 			continue
	// 		}
	// 		tasks = updatedTasks
	// 		fmt.Printf("Task %s has been done\n", params)
	// 	case "":
	// 		continue
	// 	default:
	// 		fmt.Printf("Unknown command: %s\n", cmd)
	// 	}
	// }
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standart input: ", err)
	}
}
