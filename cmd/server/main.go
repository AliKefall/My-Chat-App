package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Root -> index.html
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexPath := filepath.Join("web", "index.html")
		http.ServeFile(w, r, indexPath)
	})

	// Statik dosyalar (css, js, images vs.)
	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	fmt.Println("Server başlatıldı: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
