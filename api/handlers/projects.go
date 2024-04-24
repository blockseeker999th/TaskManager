package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blockseeker999th/TaskManager/api/auth"
	"github.com/blockseeker999th/TaskManager/db"
	"github.com/blockseeker999th/TaskManager/models"
	"github.com/blockseeker999th/TaskManager/utils"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var (
	errProjectNameRequired = errors.New("project name required")
	errIDRequired          = errors.New("project id required")
)

type ProjectService struct {
	store  db.Store
	userId string
}

func NewProjectService(s db.Store) *ProjectService {
	return &ProjectService{store: s}
}

func (s *ProjectService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/projects", auth.WithJWTAuth(s.handleCreateProject, s.store)).Methods("POST")
	r.HandleFunc("/projects/{id}", auth.WithJWTAuth(s.handleGetProject, s.store)).Methods("GET")
	r.HandleFunc("/projects/{id}", auth.WithJWTAuth(s.handleDeleteProject, s.store)).Methods("DELETE")
}

func (s *ProjectService) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	var project *models.Project
	err = json.Unmarshal(body, &project)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	err = validateProjectPayload(project)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	p, err := s.store.CreateProject(project, userID)
	fmt.Println(err)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Error creating project"})
	}

	utils.WriteJSON(w, http.StatusCreated, p)
}

func (s *ProjectService) handleGetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Project id is required"})
		return
	}

	p, err := s.store.GetProjectByID(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Error getting project by id"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, p)
}

func (s *ProjectService) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, "Project id is required")
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
