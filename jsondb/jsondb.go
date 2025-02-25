package jsondb

import (
	"encoding/json"
	"fmt"
	"go-cli-task-tracker/db"
	"io"
	"os"
	"sync"
	"syscall"
	"time"
)

// JsonTaskDB is the implementation of the TaskDB interface for a JSON file
type JsonTaskDB struct {
	// The path to the JSON file
	filePath string
	// The file
	file *os.File
	// The mutex for the file
	mu sync.RWMutex
}

// NewJsonTaskDB creates a new JsonTaskDB
func NewJsonTaskDB(filePath string) *JsonTaskDB{
	return &JsonTaskDB{
		filePath: filePath,
		mu: sync.RWMutex{},
	}
}

// Connect connects to the JSON file, creating it if it doesn't exist
func (jdb *JsonTaskDB) Connect() error {
	jdb.mu.Lock()
	defer jdb.mu.Unlock()

	file, err := os.OpenFile(jdb.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	
	jdb.file = file
	
	// Apply file level lock
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	if err != nil {
		return err
	}

	// Write an empty array to the file.
	if err = json.NewEncoder(jdb.file).Encode([]*db.Task{}); err != nil {
		return err
	}
	// Set the offset to the beginning of the file
	if _, err = jdb.file.Seek(0, 0); err != nil {
		return err
	}

	return nil
}

// Close removes the file lock and the process lock
func (jdb *JsonTaskDB) Close() error {
	jdb.mu.Lock()
	defer jdb.mu.Unlock()

	if err := syscall.Flock(int(jdb.file.Fd()), syscall.LOCK_UN); err != nil {
		return err
	}
	if err := jdb.file.Close(); err != nil {
		return err
	}
	if err := os.Remove(jdb.filePath); err != nil {
		return err
	}
	jdb.file = nil
	return nil
}

// getTasks reads the tasks from the file
func (jdb *JsonTaskDB) getTasks() ([]*db.Task, error) {
	_, err := jdb.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	tasks := make([]*db.Task, 0)
	if err = json.NewDecoder(jdb.file).Decode(&tasks); err != nil {
		if err.Error() == io.EOF.Error() {
			return []*db.Task{}, nil
		}
		return nil, err
	}

	return tasks, nil
}

// writeTasks writes the tasks to the file
func (jdb *JsonTaskDB) writeTasks(tasks []*db.Task) error {
	// Truncate the file before writing
	if err := jdb.file.Truncate(0); err != nil {
		return err
	}

	if _, err := jdb.file.Seek(0, 0); err != nil {
		return err
	}
	encoder := json.NewEncoder(jdb.file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tasks)
}

// GetTasks get all tasks from the file
func (jdb *JsonTaskDB) GetTasks() ([]*db.Task, error) {
	jdb.mu.RLock()
	defer jdb.mu.RUnlock()
	return jdb.getTasks()
}

// AddTask adds a task to the file
func (jdb *JsonTaskDB) AddTask(t *db.Task) error {
	jdb.mu.Lock()
	defer jdb.mu.Unlock()
	tasks, err := jdb.getTasks()
	if err != nil {
		return err
	}
	if len(tasks) > 0 {
		t.ID = tasks[len(tasks)-1].ID + 1
	} else {
		t.ID = uint64(time.Now().Unix())
	}
	tasks = append(tasks, t)
	return jdb.writeTasks(tasks)

}

// DeleteTask deletes a task from the file
func (jdb *JsonTaskDB) DeleteTask(taskID uint64) error {
	jdb.mu.Lock()
	defer jdb.mu.Unlock()

	tasks, err := jdb.getTasks()
	if err != nil {
		return err
	}

	for i, t := range tasks {
		if t.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:	]...)
			return jdb.writeTasks(tasks)
		}
	}
	return fmt.Errorf("task with ID %d not found", taskID)
}