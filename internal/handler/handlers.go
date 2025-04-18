package handler

import (
	"JwtAuth/internal/repo"
	"net/http"
)

type Server struct {
	db *repo.DB
}

func NewServer(db *repo.DB) *Server {
	return &Server{db: db}
}

func (s *Server) AccessMethod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		w.Write([]byte("Method - " + r.Method + " not allowed :("))
	}
}

func (s *Server) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		w.Write([]byte("Method - " + r.Method + " not allowed :("))
	}
}
