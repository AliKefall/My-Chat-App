package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/AliKefall/My-Chat-App/internal/database"
)

type config struct {
	Port string
	db   *database.Queries
}

func main() {
	mux := http.NewServeMux()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("sqlite", dbUrl)
	if err != nil {
		log.Fatalf("Error opening databse: %s", err)
	}
	dbQueries := database.New(dbConn)

	cfg := config{
		Port: ":8080:",
		db:   dbQueries,
	}

	mux.HandleFunc("POST /api/register", cfg.handleUserRegister)

}
