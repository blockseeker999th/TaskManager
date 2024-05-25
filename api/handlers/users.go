package handlers

import (
	"encoding/json"
	"errors"
	"github.com/blockseeker999th/TaskManager/api/auth"
	"github.com/blockseeker999th/TaskManager/config"
	"github.com/blockseeker999th/TaskManager/db"
	"github.com/blockseeker999th/TaskManager/models"
	"github.com/blockseeker999th/TaskManager/utils"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

var (
	errPasswordRequired  = errors.New("password required")
	errFirstNameRequired = errors.New("first name is mandatory")
	errSignUp            = errors.New("error registering a user")
)

type UserService struct {
	store db.Store
}

func NewUserService(s db.Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.ErrMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: utils.ErrReadingResponse})
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			return
		}
	}()

	var user *models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: utils.ErrInvalidRequestPayload})
	}

	hashedPassword, err := auth.HashPassword(user.Password)

	if err := validateUserRegisterPayload(user); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: err})
	}

	user.Password = hashedPassword

	u, err := s.store.CreateUser(user)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: errSignUp})
		return
	}
	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: utils.ErrCreatingSession})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, token)
}

func (s *UserService) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.ErrMethodNotAllowed)
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

	var loginData *models.LoginData
	err = json.Unmarshal(body, &loginData)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: utils.ErrInvalidRequestPayload})
		return
	}

	user, err := s.store.LoginUser(loginData)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: utils.ErrInvalidRequestPayload})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: utils.ErrUnauthorized})
		return
	}

	token, err := createAndSetAuthCookie(user.ID, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: utils.ErrCreatingSession})
		return
	}

	utils.WriteJSON(w, http.StatusOK, token)
}

func validateUserRegisterPayload(user *models.User) error {
	if user.Email == "" {
		return utils.ErrEmailRequired
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
