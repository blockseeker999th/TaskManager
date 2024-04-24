package handlers

import (
	"encoding/json"
	"errors"
	"github.com/blockseeker999th/TaskManager/api/auth"
	"github.com/blockseeker999th/TaskManager/config"
	"github.com/blockseeker999th/TaskManager/db"
	"github.com/blockseeker999th/TaskManager/models"
	"github.com/blockseeker999th/TaskManager/utils"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

var (
	errEmailRequired     = errors.New("email required")
	errPasswordRequired  = errors.New("password required")
	errFirstNameRequired = errors.New("first name is mandatory")
)

type UserService struct {
	store db.Store
}

func NewUserService(s db.Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleUserRegister).Methods("POST")
	r.HandleFunc("/users/login", s.handleUserLogin).Methods("POST")
}

func (s *UserService) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "error reading response BODY"})
		return
	}

	defer r.Body.Close()

	var user *models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
	}

	hashedPassword, err := auth.HashPassword(user.Password)

	if err := validateUserRegisterPayload(user); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	user.Password = hashedPassword

	u, err := s.store.CreateUser(user)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Error registering a user"})
		return
	}
	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Error creating session"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, token)
}

func (s *UserService) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user credentials"})
		return
	}

	defer r.Body.Close()

	var loginData *models.LoginData
	err = json.Unmarshal(body, &loginData)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	user, err := s.store.LoginUser(loginData)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid email or password"})
		return
	}

	token, err := createAndSetAuthCookie(user.ID, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Error creating session"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, token)
}

func validateUserRegisterPayload(user *models.User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	if user.FirstName == "" {
		return errFirstNameRequired
	}

	return nil
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, id)

	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
