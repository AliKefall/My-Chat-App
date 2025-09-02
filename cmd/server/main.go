package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/AliKefall/My-Chat-App/internal/ws"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// .env yükle
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// ENV değişkenlerini al
	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		log.Fatal("DB_CONN_STR environment variable not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	// Database bağlantısı
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Couldn't connect to database: ", err)
	}
	defer dbConn.Close()

	// Server oluştur
	srv := NewServer(dbConn, "8080")
	fsHandler := http.FileServer(http.Dir("./web/pages/register.html"))

	// WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	srv.Router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	// Statik sayfalar
	srv.Router.Handle("/", fsHandler)

	// API endpointleri
	srv.Router.HandleFunc("POST /api/register", srv.handleUserRegister)
	srv.Router.HandleFunc("POST /api/login", srv.handleUserLogin)

	log.Println("✅ Server is started on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":"+srv.PORT, srv.Router))
}
