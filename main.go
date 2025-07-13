package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/AliKefall/My-Chat-App/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	Port string
	db   *database.Queries
}

func main() {
	mux := http.NewServeMux()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Database url could not be loaded: %s", err)
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		log.Fatalf("Error opening databse: %s", err)
	}
	dbQueries := database.New(dbConn)

	cfg := config{
		Port: ":8080",
		db:   dbQueries,
	}

	mux.HandleFunc("POST /api/register", cfg.handleUserRegister)
	mux.HandleFunc("/ws", cfg.handlerWebsocketConn)
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)

	go handleMessages()

	log.Printf("Server is running at http://localhost%s", cfg.Port)
	err = http.ListenAndServe(cfg.Port, mux)
	if err != nil {
		log.Fatalf("Server Failed: %s", err)
	}
}
