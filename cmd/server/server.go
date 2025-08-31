package main

import (
	"database/sql"

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
