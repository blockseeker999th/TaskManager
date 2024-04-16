package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTask(test *testing.T) {
	ms := &MockStore{}
	service := NewTaskService(ms)
	test.Run("should return an error if name is empty", func(t *testing.T) {
		payload := &Task{
			Name:         "Create a new task",
			ProjectID:    1,
			AssignedToID: 42,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks", WithJWTAuth(service.handleCreateTask, ms))

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetTask(test *testing.T) {
	ms := &MockStore{}
	service := NewTaskService(ms)
	test.Run("should return an error if name is empty", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/tasks/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks", WithJWTAuth(service.handleCreateTask, ms))

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}
