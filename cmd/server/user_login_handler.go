package main

import (
	"encoding/json"
	"net/http"

	"github.com/AliKefall/My-Chat-App/crypting"
	"github.com/AliKefall/My-Chat-App/internal/database"
)

type parameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type response struct {
	User database.User
}

func (s *Server) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := s.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = crypting.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	RespondWithJSON(w, http.StatusOK, response{
		User: database.User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			Email:     user.Email,
		},
	})

}
