package server

import (
	"github.com/DAlconoiD/balance-service/storage"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	router *mux.Router
}

//New creates a server
func New() *Server {
	s := Server{router: mux.NewRouter()}
	return &s
}

//Start starts server
func (s *Server) Start(address string) error {
	return http.ListenAndServe(address, s.router)
}

//ConfigureRouter binds handles to routes
func (s *Server) ConfigureRouter(storage storage.Store) {
	s.router.HandleFunc("/hello", handleHello()).Methods("GET")
	s.router.HandleFunc("/{id:[0-9]+}", handleGetBalance(storage)).Methods("GET")
	s.router.HandleFunc("/transactions/{id:[0-9]+}", handleGetTransactions(storage)).Methods("GET")
	s.router.HandleFunc("/transfer", handleTransfer(storage)).Methods("POST")
	s.router.HandleFunc("/change-balance", handleChangeBalance(storage)).Methods("POST")
}
