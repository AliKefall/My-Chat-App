package ws

import (
	"net/http"

	"github.com/AliKefall/My-Chat-App/auth"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Lokal geliştirme için açık, prod'da sıkılaştır
	},
}

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type Client struct {
	Conn *websocket.Conn
	Send chan Message
	User string
	ID   int64
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case msg := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

// WebSocket handler
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Conn: conn,
		Send: make(chan Message),
		User: claims.Username,
		ID:   claims.UserID,
	}

	hub.Register <- client

	// Mesaj alma
	go func() {
		defer func() { hub.Unregister <- client }()
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}
			msg.User = client.User
			hub.Broadcast <- msg
		}
	}()

	// Mesaj gönderme
	go func() {
		for m := range client.Send {
			conn.WriteJSON(m)
		}
	}()
}
