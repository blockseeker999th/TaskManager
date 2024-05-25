package auth

import (
	"context"
	"fmt"
	"github.com/blockseeker999th/TaskManager/config"
	"github.com/blockseeker999th/TaskManager/models"
	"github.com/blockseeker999th/TaskManager/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := validateJWT(tokenString)

		if err != nil {
			log.Println("Failed to authenticate")
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)

		ctx := context.WithValue(r.Context(), "userID", userID)

		handlerFunc(w, r.WithContext(ctx))
	}
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteJSON(w, http.StatusUnauthorized, models.ErrorResponse{
		Error: fmt.Errorf("permission denied"),
	})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func validateJWT(t string) (*jwt.Token, error) {
	secret := config.Envs.JWTSecret
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func HashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CreateJWT(secret []byte, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", nil
	}

	return tokenString, nil
}
