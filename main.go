package main

import (
	"log"
)

func main() {
	dbConfig := Envs

	storage := ConnectToPostgreSQL(dbConfig)

	db, err := storage.InitNewPostgreSQLStorage()
	if err != nil {
		log.Fatalf("Error: %v", err)
		log.Fatal("Error to initialize PostgreSQL")
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing PostgreSQL connection: %v", err)
		}
	}()

	store := NewStore(db)
	api := NewAPIServer(":3000", store)
	api.Serve()
}
