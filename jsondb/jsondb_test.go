package jsondb

import (
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
	
	db := NewJsonTaskDB(tempFile.Name())
	if err := db.Connect(); err != nil {
		t.Fatalf("Failed to connect to db: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		t.Fatalf("Failed to close db: %s", err.Error())
	}
}

func TestJsonTaskDB_GetTasks(t *testing.T) {
	path := t.TempDir()
	tempFile, err := os.CreateTemp(path, "*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err.Error())
	}
	defer os.Remove(tempFile.Name())	
	
	db := NewJsonTaskDB(tempFile.Name())
	if err := db.Connect(); err != nil {
		t.Fatalf("Failed to connect to db: %s", err.Error())
	}

	tasks := []db.Task{
		{
			ID: 1,
			Title: "Task 1",
			IsCompleted: false,
		},
		{
			ID: 2,
			Title: "Task 2",
			IsCompleted: true,
		},
		{
			ID: 3,
			Title: "Task 3",
			IsCompleted: false,
		},
	}
	
	for _, task := range tasks {
		if err := db.AddTask(&task); err != nil {
			t.Fatalf("Failed to add task: %s", err.Error())
		}
	}

	gotTasks, err := db.GetTasks()
	if err != nil {
		t.Fatalf("Failed to get tasks: %s", err.Error())	
	}

	if len(gotTasks) != len(tasks) {
		t.Fatalf("Expected %d tasks, got %d", len(tasks), len(gotTasks))
	}

	for i, task := range tasks {
		if gotTasks[i].ID != tasks[i].ID {
			t.Fatalf("Expected task %d to have ID %d, got %d", i, tasks[i].ID, gotTasks[i].ID)		
		}
		if gotTasks[i].Title != tasks[i].Title {
			t.Fatalf("Expected task %d to have title %s, got %s", i, tasks[i].Title, gotTasks[i].Title)
		}
		if gotTasks[i].IsCompleted != tasks[i].IsCompleted {
			t.Fatalf("Expected task %d to have IsCompleted %t, got %t", i, tasks[i].IsCompleted, gotTasks[i].IsCompleted)	
		}
	}
	
	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close db: %s", err.Error())
	}
}