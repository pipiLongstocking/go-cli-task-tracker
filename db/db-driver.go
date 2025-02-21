package db

// TaskDB defines the interface for the database driver
type TaskDB interface {
	// Connect to the database
	Connect() error
	// Close the connection to the database	
	Close() error
	// AddTask adds a task to the database
	AddTask(task *Task) error
	// GetTasks gets all tasks from the database
	GetTasks() ([]*Task, error)
	// DeleteTask deletes a task from the database
	DeleteTask(taskID uint64) error
}
