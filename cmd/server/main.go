package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/AliKefall/My-Chat-App/internal/ws"
	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DB_CONN_STR")
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Couldn't connect to database: ", err)
	}
	defer dbConn.Close()

	srv := NewServer(dbConn, "8080")

	// Websocket hub
	hub := ws.NewHub()
	go hub.Run()

	// Routes for frontend pages
	srv.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/pages/index.html")
	})

	srv.Router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/pages/login.html")
	})

	srv.Router.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/pages/register.html")
	})

	srv.Router.Get("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/pages/chat.html")
	})

	// Static files (css, js, images)
	srv.Router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/public"))))

	// API endpoints
	srv.Router.HandleFunc("POST /api/register", srv.handleUserRegister)
	srv.Router.HandleFunc("POST /api/login", srv.handleUserLogin)

	// WebSocket endpoint
	srv.Router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	// Start server
	log.Println("Server is started on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":"+srv.PORT, srv.Router))
}
