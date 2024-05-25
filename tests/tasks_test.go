package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/blockseeker999th/TaskManager/api/handlers"
	"github.com/blockseeker999th/TaskManager/models"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*func mockTokenValidation() *jwt.Token {
	claims := jwt.MapClaims{
		"userID": "42",                                  // Стандартне значення ідентифікатора користувача для тестування
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Стандартний термін дії токена для тестування
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}*/

func TestCreateTask(test *testing.T) {
	ms := &MockStore{}
	service := handlers.NewTaskService(ms)

	payload := &models.Task{
		Name:         "Create a new task",
		Status:       "TODO",
		ProjectID:    1,
		AssignedToID: 42,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		test.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
	if err != nil {
		test.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/tasks", service.HandleCreateTask)

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		test.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
	}

	fmt.Println("Response body:", rr.Body.String())
}

func TestGetTask(test *testing.T) {
	ms := &MockStore{}
	service := handlers.NewTaskService(ms)

	router := mux.NewRouter()
	router.HandleFunc("/tasks/42", service.HandleGetTask)

	req, err := http.NewRequest(http.MethodGet, "/tasks/42", nil)
	if err != nil {
		test.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		test.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}
