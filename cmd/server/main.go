package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
)

func main() {
	connStr := os.Getenv("DB_CONN_STR")
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Couldn't connect to database: ", err)
	}
	defer dbConn.Close()

	srv := NewServer(dbConn, "8080")
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	srv.Router.Handle("/app/", fsHandler)

	srv.Router.HandleFunc("POST /api/register", srv.handleUserRegister)
	srv.Router.HandleFunc("POST /api/login", srv.handleUserLogin)
	log.Println("Server is started on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":"+srv.PORT, srv.Router))
}
