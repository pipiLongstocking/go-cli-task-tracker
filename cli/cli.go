package cli

import (
	"fmt"
	"go-cli-task-tracker/db"
	"os"
	"strconv"
)

func getTaskStatus(done bool) string {
	if done {
		return "x"
	}
	return " "
}

func Run(taskDB db.TaskDB) {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ttracker <command>")
		fmt.Println("Commands:")	
		fmt.Println("  add - Add a new task")
		fmt.Println("  list - List all tasks")
		fmt.Println("  delete <id> - Delete a task")
		fmt.Println("  complete <id> - Mark a task as completed")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ttracker add <title>")
			os.Exit(1)
		}
		title := os.Args[2]
		err := taskDB.AddTask(&db.Task{Title: title})
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Task added successfully")
		
	case "list":
		tasks, err := taskDB.GetTasks()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		for _, t := range tasks {
			fmt.Printf("%d.)\t%s\t[%s]\n", t.ID, t.Title, getTaskStatus(t.IsCompleted))
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ttracker delete <id>")
			os.Exit(1)
		}
		id, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println("Invalid task ID")
			os.Exit(1)
		}
		err = taskDB.DeleteTask(id)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ttracker complete <id>")
			os.Exit(1)
		}
		id, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println("Invalid task ID")
			os.Exit(1)
		}
		err = taskDB.CompleteTask(id)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Invalid command")
		fmt.Println("Usage: ttracker <command>")
		fmt.Println("Commands:")	
		fmt.Println("  add - Add a new task")
		fmt.Println("  list - List all tasks")
		fmt.Println("  delete <id> - Delete a task")
		fmt.Println("  complete <id> - Mark a task as completed")
		os.Exit(1)
	}
}