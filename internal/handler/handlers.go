package handler

import (
	"JwtAuth/internal/domain"
	"JwtAuth/internal/repo"
	"JwtAuth/internal/token"
	"encoding/json"
	"log"
	"net/http"
	"time"
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
	w.Header().Set("Content-Type", "application/json")

	var body struct {
		GuId string `json:"gu_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	accessToken, jti, err := token.GenerateAccessToken(body.GuId, r.RemoteAddr, time.Minute*5)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, hashedRefreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	err = s.db.StoreRefreshToken(body.GuId, hashedRefreshToken, jti, r.RemoteAddr)
	if err != nil {
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		w.Write([]byte("Method - " + r.Method + " not allowed :("))
	}
	w.Header().Set("Content-Type", "application/json")

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	refreshToken, err := s.db.FindRefreshToken(body.RefreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	if refreshToken.IPAddress != r.RemoteAddr {
		log.Println("АААААА, Другой айпи")
		// Отправка сообзения по Email
	}

	newAccessToken, newJTI, err := token.GenerateAccessToken(refreshToken.UserGuid, r.RemoteAddr, time.Minute*5)
	if err != nil {
		http.Error(w, "could not create access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, newHash, err := token.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "could not create refresh token", http.StatusInternalServerError)
		return
	}

	err = s.db.UpdateRefreshToken(refreshToken.ID, refreshToken.UserGuid, newHash, newJTI, r.RemoteAddr)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.Token{AccessToken: newAccessToken, RefreshToken: newRefreshToken})

}
