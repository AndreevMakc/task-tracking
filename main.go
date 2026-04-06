package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"task-tracking/internal/domain"
	"task-tracking/internal/repository"
	"task-tracking/internal/service"
	"task-tracking/internal/storage/file"
)

type App struct {
	namespace           domain.Namespace
	namespaceRepository repository.NamespaceRepository
}

func newApp() (*App, error) {
	app := App{}
	namespaceRepo, err := file.NewNamespaceRepository(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to init repository: %w", err)
	}
	app.namespaceRepository = namespaceRepo
	var currentNamespace string
	if len(os.Args) > 1 {
		currentNamespace = os.Args[1]
	}
	if currentNamespace == "" {
		currentNamespace = file.DefaultNamespaceName
	}
	namespace, err := app.namespaceRepository.FindByName(currentNamespace)
	if err != nil {
		return nil, fmt.Errorf("namespace not found: %w", err)
	}
	if namespace != nil {
		app.namespace = *namespace
	} else {
		app.namespace = domain.Namespace{Name: currentNamespace}
		if err := app.namespaceRepository.Save(app.namespace); err != nil {
			return nil, fmt.Errorf("failed to create namespace: %w", err)
		}
	}
	return &app, nil
}

func humanId(namespace string, seqId int) string {
	return fmt.Sprintf("%s-%d", namespace, seqId)
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
			if err := app.namespaceRepository.Save(app.namespace); err != nil {
				fmt.Fprintf(os.Stderr, "File save error: %s\n", err)
			}
			fmt.Println("Goodbye")
			return
		case "add":
			task, err := service.CreateTask(params, &app.namespace)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Task add error: %s\n", err)
				continue
			}
			app.namespace.Tasks = append(app.namespace.Tasks, task)
		case "list":
			if len(app.namespace.Tasks) == 0 {
				fmt.Println("No tasks yet")
				continue
			}
			for _, v := range app.namespace.Tasks {
				fmt.Printf("[%s] %s - %t\n", humanId(app.namespace.Name, v.SeqId), v.Title, v.IsDone)
			}
		case "done":
			if params == "" {
				fmt.Fprintf(os.Stderr, "task id is required")
				continue
			}
			if err := service.TaskDone(params, &app.namespace); err != nil {
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
