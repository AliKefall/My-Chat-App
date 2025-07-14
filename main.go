package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AliKefall/My-Chat-App/internal/database"
	"github.com/joho/godotenv"
	"github.com/tursodatabase/libsql-client-go/libsql"
)

type config struct {
	Port string
	db   *database.Queries
}

func main() {
	// .env dosyasını yükle (Geliştirme ortamı için)
	err := godotenv.Load()
	if err != nil {
		log.Println(".env dosyası yüklenemedi, devam ediliyor (production ortamında bu normaldir)")
	}

	// Ortam değişkenlerinden DB bağlantı adresini al
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("HATA: DB_URL ortam değişkeni ayarlanmalı")
	}

	// LibSQL (Turso) bağlantısı oluştur
	dbConn, err := libsql.Open("libsql", dbURL)
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}
	defer dbConn.Close()

	// Sorgular için wrapper oluştur
	dbQueries := database.New(dbConn)

	// Yapılandırma nesnesi oluştur
	cfg := config{
		Port: ":8080",
		db:   dbQueries,
	}

	// Router oluştur
	mux := http.NewServeMux()

	// API ve WebSocket uç noktaları
	mux.HandleFunc("POST /api/register", cfg.handleUserRegister)
	mux.HandleFunc("/wss", cfg.handlerWebsocketConn)

	// Statik frontend dosyaları
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)

	// Mesajları yöneten goroutine başlat
	go handleMessages()

	log.Printf("Sunucu çalışıyor: http://localhost%s", cfg.Port)
	err = http.ListenAndServe(cfg.Port, mux)
	if err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}
