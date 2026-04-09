package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"task-tracking/internal/domain"
	"task-tracking/internal/service"
	"task-tracking/internal/storage/filestorage"
)

type App struct {
	taskService *service.TaskService
}

func getNamespace() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return domain.DefaultNamespace
}

func newApp() (*App, error) {
	fileNamespaceRepo, err := filestorage.NewFileNamespaceRepository("")
	if err != nil {
		return nil, fmt.Errorf("init app failed: %w", err)
	}
	taskService := service.NewTaskService(fileNamespaceRepo)
	app := App{taskService: taskService}
	return &app, nil
}

func handleInputCmd(text string) (cmd string, taskTitle string) {
	sliceText := strings.SplitN(text, " ", 2)
	cmd = sliceText[0]
	if len(sliceText) > 1 {
		taskTitle = sliceText[1]
	}
	return cmd, taskTitle
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
			fmt.Println("Goodbye")
			return
		case "add":
			if _, err := app.taskService.CreateTask(params, getNamespace()); err != nil {
				fmt.Fprintf(os.Stderr, "Task add error: %s\n", err)
				continue
			}
		case "list":
			for _, v := range app.taskService.ShowList(getNamespace()) {
				fmt.Println(v)
			}
		case "done":
			if params == "" {
				fmt.Fprintf(os.Stderr, "task id is required")
				continue
			}
			if err := app.taskService.TaskDone(params, getNamespace()); err != nil {
				fmt.Fprintf(os.Stderr, "Task done failed: %s\n", err)
				continue
			}
		case "":
			continue
		default:
			fmt.Printf("unknown command: %s\n", cmd)
		}
	}
}
