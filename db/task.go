package db

// Task defines the structure of the task object
type Task struct {
	ID          uint64	`json:"id"`
	Title       string	`json:"title"`
	// IsCompleted denotes if the task has been completed
	IsCompleted bool	`json:"done"`
}