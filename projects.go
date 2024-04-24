package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var (
	errProjectNameRequired = errors.New("project name required")
	errIDRequired          = errors.New("project id required")
)

type ProjectService struct {
	store  Store
	userId string
}

func NewProjectService(s Store) *ProjectService {
	return &ProjectService{store: s}
}

func (s *ProjectService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/projects", WithJWTAuth(s.handleCreateProject, s.store)).Methods("POST")
	r.HandleFunc("/projects/{id}", WithJWTAuth(s.handleGetProject, s.store)).Methods("GET")
	r.HandleFunc("/projects/{id}", WithJWTAuth(s.handleDeleteProject, s.store)).Methods("DELETE")
}

func (s *ProjectService) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	var project *Project
	err = json.Unmarshal(body, &project)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	err = validateProjectPayload(project)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	p, err := s.store.CreateProject(project, userID)
	fmt.Println(err)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating project"})
	}

	WriteJSON(w, http.StatusCreated, p)
}

func (s *ProjectService) handleGetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Project id is required"})
		return
	}

	p, err := s.store.GetProjectByID(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error getting project by id"})
		return
	}

	WriteJSON(w, http.StatusOK, p)
}

func (s *ProjectService) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		WriteJSON(w, http.StatusBadRequest, "Project id is required")
		return
	}

	err := s.store.DeleteProjectByID(id, userID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSON(w, http.StatusAccepted, id)
}

func validateProjectPayload(p *Project) error {
	if p.Name == "" {
		return errProjectNameRequired
	}

	return nil
}
