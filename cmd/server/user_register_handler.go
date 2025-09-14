package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AliKefall/My-Chat-App/auth"
	"github.com/AliKefall/My-Chat-App/crypting"
	"github.com/AliKefall/My-Chat-App/internal/database"
)

type RegisteredRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisteredRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	passwordHash, err := crypting.HashPassword(req.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't hash the password", err)
	}

	user, err := s.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
	})

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Println("Register error : ", err)
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
