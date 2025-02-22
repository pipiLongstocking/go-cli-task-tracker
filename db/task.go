package db

type Task struct {
	ID          uint64	`json:"id"`
	Title       string	`json:"title"`
	Description string	`json:"description"`
	Done        bool	`json:"done"`
}