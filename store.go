package main

import (
	"database/sql"
)

type Store interface {
	CreateUser(u *User) (*User, error)
	GetUserByID(id string) (*User, error)
	CreateTask(t *Task) (*Task, error)
	GetTask(id string) (*Task, error)
	LoginUser(data *LoginData) (*LoginData, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateTask(t *Task) (*Task, error) {
	err := s.db.QueryRow("INSERT INTO tasks (name, status, projectId, assignedToID) VALUES ($1, $2, $3, $4) RETURNING id",
		t.Name, t.Status, t.ProjectID, t.AssignedToID).Scan(&t.ID)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task
	err := s.db.QueryRow("SELECT id, name, status, projectid, assignedtoid, createdat FROM tasks WHERE id = $1", id).Scan(
		&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID, &t.CreatedAt)
	return &t, err
}

func (s *Storage) GetUserByID(id string) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT id, firstname, lastname, password, createdat FROM users WHERE id = $1", id).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Password, &u.CreatedAt)
	return &u, err
}

func (s *Storage) CreateUser(u *User) (*User, error) {
	err := s.db.QueryRow("INSERT INTO users (email, firstname, lastname, password) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Email, u.FirstName, u.LastName, u.Password).Scan(&u.ID)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Storage) LoginUser(data *LoginData) (*LoginData, error) {
	var l LoginData
	err := s.db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", data.Email).Scan(&l.ID, &l.Email, &l.Password)

	return &l, err
}
