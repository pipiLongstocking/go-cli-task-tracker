package jsondb

import (
	"os"
	"sync"
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
func (db *JsonTaskDB) Connect() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	file, err := os.OpenFile(db.filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	
	db.file = file
	
}
