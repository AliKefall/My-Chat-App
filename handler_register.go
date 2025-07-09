package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AliKefall/My-Chat-App/internal/auth"
	"github.com/AliKefall/My-Chat-App/internal/database"
)

func (cfg *config) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request data!")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "All fields must be filled.")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password could not be hashed!")
		return
	}

	user := database.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	err = cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "User has been saved succesfully",
	})
}
