package tests

import "github.com/blockseeker999th/TaskManager/models"

type MockStore struct {
}

func (m *MockStore) CreateUser(u *models.User) (*models.User, error) {
	return &models.User{}, nil
}

func (m *MockStore) CreateTask(t *models.Task) (*models.Task, error) {
	return &models.Task{}, nil
}

func (m *MockStore) GetTask(id string) (*models.Task, error) {
	return &models.Task{}, nil
}

func (m *MockStore) GetUserByID(id string) (*models.User, error) {
	return &models.User{}, nil
}
