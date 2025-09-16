package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AliKefall/My-Chat-App/internal/ws"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	if err := godotenv.Load(filepath.Join(".", ".env")); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		log.Fatal("❌ DB_CONN_STR environment variable not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET environment variable not set")
	}

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Couldn't connect to database: ", err)
	}
	defer dbConn.Close()

	srv := NewServer(dbConn, "8080")

	hub := ws.NewHub()
	go hub.Run()

	srv.Router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	fileServerPages := http.FileServer(http.Dir("./web/pages"))
	srv.Router.Handle("/*", fileServerPages)

	fileServerStatic := http.FileServer(http.Dir("./web/static"))
	srv.Router.Handle("/static/*", http.StripPrefix("/static/", fileServerStatic))

	srv.Router.HandleFunc("POST /api/register", srv.handleUserRegister)

	srv.Router.Handle("POST /api/login", http.HandlerFunc(srv.handleUserLogin))

	log.Println("✅ Server is running at: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":"+srv.PORT, srv.Router))
}
