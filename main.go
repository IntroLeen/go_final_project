package main

import (
	"log"
	"os"

	mydb "go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := mydb.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
