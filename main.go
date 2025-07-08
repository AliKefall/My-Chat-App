package main

import (
	"net/http"

	"github.com/AliKefall/My-Chat-App/internal/config"
)

type cfg struct {
	user config.User
}

func main() {
	mux := http.NewServeMux()

}
