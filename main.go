package main

import (
	"fmt"
	"go-cli-task-tracker/cli"
	"go-cli-task-tracker/jsondb"
	"os"
)

func main() {
	taskDB := jsondb.NewJsonTaskDB("tasks.json")
	if err := taskDB.Connect(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer taskDB.Close()

	cli.Run(taskDB)
}