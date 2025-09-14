package main

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/AliKefall/My-Chat-App/auth"
	"github.com/AliKefall/My-Chat-App/internal/database"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	DB     *database.Queries
	Router *chi.Mux
	PORT   string
}

func NewServer(dbConn *sql.DB, port string) *Server {
	s := &Server{
		DB:     database.New(dbConn),
		PORT:   port,
		Router: chi.NewRouter(),
	}

	return s
}

// This part is totally unnecessary for now but will be used in the future. It is a middleware that checks for the presence of a valid JWT token in the Authorization header of incoming requests. If the token is valid, it adds the user information to the request context and calls the next handler. If the token is missing or invalid, it responds with an unauthorized error.
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
