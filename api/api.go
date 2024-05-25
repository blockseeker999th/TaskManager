package api

import (
	"github.com/blockseeker999th/TaskManager/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	addr  string
	store db.Store
}

func NewAPIServer(addr string, store db.Store) *APIServer {
	return &APIServer{addr: addr, store: store}
}

func (s *APIServer) Serve() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	
	RegisterRoutes(subrouter, s.store)

	log.Println("Starting the API server at: ", s.addr)

	log.Fatal(http.ListenAndServe(s.addr, subrouter))
}
