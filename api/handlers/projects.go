package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blockseeker999th/TaskManager/db"
	"github.com/blockseeker999th/TaskManager/models"
	"github.com/blockseeker999th/TaskManager/utils"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var (
	errProjectNameRequired = errors.New("project name required")
	errCreatingProject     = errors.New("error creating a project")
	errProjectNotFound     = errors.New("not found project with such id")
)

type ProjectService struct {
	store  db.Store
	userId string
}

func NewProjectService(s db.Store) *ProjectService {
	return &ProjectService{store: s}
}

func (s *ProjectService) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, models.ErrorResponse{Error: utils.ErrMethodNotAllowed})
		return
	}

	userID := r.Context().Value("userID").(string)

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

	var project *models.Project
	err = json.Unmarshal(body, &project)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: utils.ErrInvalidRequestPayload})
		return
	}

	err = validateProjectPayload(project)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: err})
		return
	}

	p, err := s.store.CreateProject(project, userID)
	fmt.Println(err)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: errCreatingProject})
	}

	utils.WriteJSON(w, http.StatusCreated, p)
}

func (s *ProjectService) HandleGetProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, models.ErrorResponse{Error: utils.ErrMethodNotAllowed})
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: utils.ErrProjectIDRequired})
		return
	}

	p, err := s.store.GetProjectByID(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: errProjectNotFound})
		return
	}

	utils.WriteJSON(w, http.StatusOK, p)
}

func (s *ProjectService) HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, models.ErrorResponse{Error: utils.ErrMethodNotAllowed})
		return
	}

	userID := r.Context().Value("userID").(string)

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrProjectIDRequired)
		return
	}

	err := s.store.DeleteProjectByID(id, userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, id)
}

func validateProjectPayload(p *models.Project) error {
	if p.Name == "" {
		return errProjectNameRequired
	}

	return nil
}
