package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type PostgreSQLStorage struct {
	db *sql.DB
}

func ConnectToPostgreSQL(config PostgreSQLConfig) *PostgreSQLStorage {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatalf("error connecting to PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("error pinging PostgreSQL: %v", err)
	}

	return &PostgreSQLStorage{db: db}
}

func (s *PostgreSQLStorage) InitNewPostgreSQLStorage() (*sql.DB, error) {
	if err := s.createProjectsTable(); err != nil {
		return nil, err
	}

	if err := s.createUsersTable(); err != nil {
		return nil, err
	}

	if err := s.createTasksTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *PostgreSQLStorage) createProjectsTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    assignedToID INT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (assignedToID) REFERENCES users(id)
)`)

	return err
}

func (s *PostgreSQLStorage) createTasksTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'TODO',
    projectId INT NOT NULL,
	assignedToID INT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT valid_status CHECK (status IN ('TODO', 'IN_PROGRESS', 'IN_TESTING', 'DONE')),
    
    FOREIGN KEY (assignedToID) REFERENCES users(id),
    FOREIGN KEY (projectId) REFERENCES projects(id)
)`)

	return err
}

func (s *PostgreSQLStorage) createUsersTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE (email)
)`)
	return err
}
