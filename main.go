package main

import (
	api2 "github.com/blockseeker999th/TaskManager/api"
	"github.com/blockseeker999th/TaskManager/config"
	db2 "github.com/blockseeker999th/TaskManager/db"
	"log"
)

func main() {
	dbConfig := config.Envs

	storage := db2.ConnectToPostgreSQL(dbConfig)

	db, err := storage.InitNewPostgreSQLStorage()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing PostgreSQL connection: %v", err)
		}
	}()

	store := db2.NewStore(db)
	api := api2.NewAPIServer(":3000", store)
	api.Serve()
}
