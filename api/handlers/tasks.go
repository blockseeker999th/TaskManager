package handlers

import (
	"encoding/json"
	"errors"
	"github.com/blockseeker999th/TaskManager/db"
	"github.com/blockseeker999th/TaskManager/models"
	"github.com/blockseeker999th/TaskManager/utils"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var (
	errTaskNameRequired = errors.New("task name required")
	errTaskIdRequired   = errors.New("task id required")
	errCreatingTask     = errors.New("error creating a task")
	errNotFoundTask     = errors.New("not found task with such id")
)

type TaskService struct {
	store db.Store
}

func NewTaskService(s db.Store) *TaskService {
	return &TaskService{store: s}
}

func (s *TaskService) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, models.ErrorResponse{Error: utils.ErrMethodNotAllowed})
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: utils.ErrUnauthorized})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: utils.ErrInvalidRequestPayload})
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
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: utils.ErrInvalidRequestPayload})
		return
	}

	if err := validateTaskPayload(task); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: err})
		return
	}

	t, err := s.store.CreateTask(task, userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: errCreatingTask})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, t)
}

func (s *TaskService) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, models.ErrorResponse{Error: utils.ErrMethodNotAllowed})
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: errTaskIdRequired})
		return
	}

	t, err := s.store.GetTask(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: errNotFoundTask})
		return
	}

	utils.WriteJSON(w, http.StatusOK, t)
}

func validateTaskPayload(t *models.Task) error {
	if t.Name == "" {
		return errTaskNameRequired
	}

	if t.ProjectID == 0 {
		return utils.ErrProjectIDRequired
	}

	return nil
}
