package handler

import (
	"JwtAuth/internal/repo"
	"JwtAuth/internal/token"
	"encoding/json"
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

	accessToken, err := token.GenerateAccessToken(body.GuId, r.RemoteAddr, time.Minute*5)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, hashedRefreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	err = s.db.StoreRefreshToken(body.GuId, hashedRefreshToken, r.RemoteAddr, "")
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

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	user, err := s.db.FindUserByRefreshToken(body.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to find user", http.StatusNotFound)
		return
	}

	if user.Ip != r.RemoteAddr {
		// Отправляем на почту уведомление, что вйпишник другой
	}

	newAccessToken, err := token.GenerateAccessToken(user.Guid, r.RemoteAddr, time.Minute*5)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, newHashedRefreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	err = s.db.UpdateRefreshToken(user.Guid, newHashedRefreshToken, r.RemoteAddr)
	if err != nil {
		http.Error(w, "Failed to update refresh token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
