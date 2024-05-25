package api

import (
	"github.com/blockseeker999th/TaskManager/api/auth"
	"github.com/blockseeker999th/TaskManager/api/handlers"
	"github.com/blockseeker999th/TaskManager/db"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, store db.Store) {
	ps := handlers.NewProjectService(store)
	ts := handlers.NewTaskService(store)
	us := handlers.NewUserService(store)

	r.HandleFunc("/projects", auth.WithJWTAuth(ps.HandleCreateProject)).Methods("POST")
	r.HandleFunc("/projects/{id}", auth.WithJWTAuth(ps.HandleGetProject)).Methods("GET")
	r.HandleFunc("/projects/{id}", auth.WithJWTAuth(ps.HandleDeleteProject)).Methods("DELETE")

	r.HandleFunc("/tasks", auth.WithJWTAuth(ts.HandleCreateTask)).Methods("POST")
	r.HandleFunc("/tasks/{id}", auth.WithJWTAuth(ts.HandleGetTask)).Methods("GET")

	r.HandleFunc("/users/register", us.HandleUserRegister).Methods("POST")
	r.HandleFunc("/users/login", us.HandleUserLogin).Methods("POST")
}
