package jsondb

import (
	"go-cli-task-tracker/db"
	"os"
	"testing"
)

func TestJsonTaskDB_ConnectAndClose(t *testing.T) {
	path := t.TempDir()
	tempFile, err := os.CreateTemp(path, "*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err.Error())
	}
	defer os.Remove(tempFile.Name())
	
	tdb := NewJsonTaskDB(tempFile.Name())
	if err := tdb.Connect(); err != nil {
		t.Fatalf("Failed to connect to db: %s", err.Error())
	}

	if err = tdb.Close(); err != nil {
		t.Fatalf("Failed to close db: %s", err.Error())
	}
}

func TestJsonTaskDB_TaskMethods(t *testing.T) {
	// Setup
	path := t.TempDir()
	tempFile, err := os.CreateTemp(path, "*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err.Error())
	}
	defer os.Remove(tempFile.Name())	
	
	// Create a new task database
	tdb := NewJsonTaskDB(tempFile.Name())
	if err := tdb.Connect(); err != nil {
		t.Fatalf("Failed to connect to db: %s", err.Error())
	}

	tasks := []db.Task{
		{
			Title: "Task 1",
			IsCompleted: false,
		},
		{
			Title: "Task 2",
			IsCompleted: true,
		},
		{
			Title: "Task 3",
			IsCompleted: false,
		},
	}
	
	// Test Add tasks
	for _, task := range tasks {
		if err := tdb.AddTask(&task); err != nil {
			t.Fatalf("Failed to add task: %s", err.Error())
		}
	}


	// Test Get all tasks
	gotTasks, err := tdb.GetTasks()
	if err != nil {
		t.Fatalf("Failed to get tasks: %s", err.Error())	
	}

	if len(gotTasks) != len(tasks) {
		t.Fatalf("Expected %d tasks, got %d", len(tasks), len(gotTasks))
	}

	for i := range tasks {
		if gotTasks[i].Title != tasks[i].Title {
			t.Fatalf("Expected task %d to have title %s, got %s", i, tasks[i].Title, gotTasks[i].Title)
		}
		if gotTasks[i].IsCompleted != tasks[i].IsCompleted {
			t.Fatalf("Expected task %d to have IsCompleted %t, got %t", i, tasks[i].IsCompleted, gotTasks[i].IsCompleted)	
		}
	}

	// Test Delete task
	err = tdb.DeleteTask(4)
	if err == nil {
		t.Fatalf("Expected error for deleting non-existent task, got nil")
	} else {
		if err.Error() != "task with ID 4 not found" {	
			t.Fatalf("Expected error message to be 'task with ID 4 not found', got %s", err.Error())
		}
	}


	// Test Delete task with valid ID
	for _, task := range gotTasks {
	err = tdb.DeleteTask(task.ID)
	if err != nil {
			t.Fatalf("Failed to delete task: %s", err.Error())
		}
	}

	// Test Get all tasks after deletion
	gotTasks, err = tdb.GetTasks()
	if err != nil {
		t.Fatalf("Failed to get tasks: %s", err.Error())
	}
	
	if len(gotTasks) != 0 {
		t.Fatalf("Expected 0 tasks, got %d", len(gotTasks))
	}	
	
	if err := tdb.Close(); err != nil {
		t.Fatalf("Failed to close db: %s", err.Error())
	}
}