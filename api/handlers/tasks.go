package handlers

import (
	"encoding/json"
	"errors"
	"github.com/blockseeker999th/TaskManager/api/auth"
	"github.com/blockseeker999th/TaskManager/db"
	"github.com/blockseeker999th/TaskManager/models"
	"github.com/blockseeker999th/TaskManager/utils"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var (
	errTaskNameRequired  = errors.New("task name required")
	errProjectIDRequired = errors.New("project id required")
	errUserIDRequired    = errors.New("user id required")
)

type TaskService struct {
	store db.Store
}

func NewTaskService(s db.Store) *TaskService {
	return &TaskService{store: s}
}

func (s *TaskService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", auth.WithJWTAuth(s.handleCreateTask, s.store)).Methods("POST")
	r.HandleFunc("/tasks/{id}", auth.WithJWTAuth(s.handleGetTask, s.store)).Methods("GET")
}

func (s *TaskService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			return
		}
	}()

	var task *models.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateTaskPayload(task); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	t, err := s.store.CreateTask(task, userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Error creating task"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, t)
}

func (s *TaskService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "id is required"})
		return
	}

	t, err := s.store.GetTask(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "no such a task with this id"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, t)
}

func validateTaskPayload(t *models.Task) error {
	if t.Name == "" {
		return errTaskNameRequired
	}

	if t.ProjectID == 0 {
		return errProjectIDRequired
	}

	return nil
}
