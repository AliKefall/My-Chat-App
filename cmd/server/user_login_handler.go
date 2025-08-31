package main

import (
	"encoding/json"
	"net/http"

	"github.com/AliKefall/My-Chat-App/auth"
)

func (s *Server) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user, err := s.DB.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateJWT(user.Username, int64(user.ID))
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
