package db

import (
	"database/sql"
	"errors"
	"github.com/blockseeker999th/TaskManager/models"
)

type Store interface {
	CreateProject(p *models.Project, userID string) (*models.Project, error)
	GetProjectByID(id string) (*models.Project, error)
	DeleteProjectByID(id string, userID string) error
	CreateTask(t *models.Task, userID string) (*models.Task, error)
	GetTask(id string) (*models.Task, error)
	CreateUser(u *models.User) (*models.User, error)
	LoginUser(data *models.LoginData) (*models.LoginData, error)
	/*GetUserByID(id string) (*User, error)*/
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateProject(p *models.Project, userID string) (*models.Project, error) {
	err := s.db.QueryRow("INSERT INTO projects (name, createdat, assignedtoid) VALUES ($1, $2, $3) RETURNING id",
		p.Name, p.CreatedAt, userID).Scan(&p.ID)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Storage) GetProjectByID(id string) (*models.Project, error) {
	var p models.Project
	err := s.db.QueryRow("SELECT id, name, createdat FROM projects WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *Storage) DeleteProjectByID(id string, userID string) error {
	res, err := s.db.Exec("DELETE FROM projects WHERE id = $1 AND assignedtoid = $2", id, userID)

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("unauthorized access")
	}

	return nil
}

func (s *Storage) CreateTask(t *models.Task, userID string) (*models.Task, error) {
	err := s.db.QueryRow("INSERT INTO tasks (name, status, projectid, assignedtoid) SELECT $1, $2, $3, projects.assignedtoid FROM projects WHERE projects.id = $3 AND projects.assignedtoid = $4 RETURNING id",
		t.Name, t.Status, t.ProjectID, userID).Scan(&t.ID)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Storage) GetTask(id string) (*models.Task, error) {
	var t models.Task
	err := s.db.QueryRow("SELECT id, name, status, projectid, assignedtoid, createdat FROM tasks WHERE id = $1", id).Scan(
		&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &t, err
}

func (s *Storage) CreateUser(u *models.User) (*models.User, error) {
	err := s.db.QueryRow("INSERT INTO users (email, firstname, lastname, password) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Email, u.FirstName, u.LastName, u.Password).Scan(&u.ID)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Storage) LoginUser(data *models.LoginData) (*models.LoginData, error) {
	var l models.LoginData
	err := s.db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", data.Email).Scan(&l.ID, &l.Email, &l.Password)

	return &l, err
}

/*func (s *Storage) GetUserByID(id string) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT id, firstname, lastname, password, createdat FROM users WHERE id = $1", id).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Password, &u.CreatedAt)
	return &u, err
}*/
