package main

import (
	"encoding/json"
	"net/http"

	"github.com/AliKefall/My-Chat-App/internal/auth"
	"github.com/AliKefall/My-Chat-App/internal/config"
)

func handleUserRegister(w http.ResponseWriter, r *http.Request) {
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

	user := config.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

}
