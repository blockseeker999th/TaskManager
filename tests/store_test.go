package tests

import "github.com/blockseeker999th/TaskManager/models"

type MockStore struct {
}

func (m *MockStore) CreateProject(p *models.Project, userID string) (*models.Project, error) {
	return &models.Project{}, nil
}

func (m *MockStore) GetProjectByID(id string) (*models.Project, error) {
	return &models.Project{}, nil
}

func (m *MockStore) DeleteProjectByID(id string, userID string) error {
	return nil
}

func (m *MockStore) CreateTask(t *models.Task, userID string) (*models.Task, error) {
	return &models.Task{}, nil
}

func (m *MockStore) GetTask(id string) (*models.Task, error) {
	return &models.Task{}, nil
}

func (m *MockStore) CreateUser(u *models.User) (*models.User, error) {
	return &models.User{}, nil
}

func (m *MockStore) LoginUser(data *models.LoginData) (*models.LoginData, error) {
	return &models.LoginData{}, nil
}
