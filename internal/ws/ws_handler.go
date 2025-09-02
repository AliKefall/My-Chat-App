package ws

import (
	"log"
	"net/http"

	"github.com/AliKefall/My-Chat-App/auth"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // geliştirme için açık, prod ortamında domain check yap
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

// ✅ readPump: sadece istemciden okuma yapar
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			log.Println("read error:", err)
			break
		}
		msg.User = c.User
		hub.Broadcast <- msg
	}
}

// ✅ writePump: sadece client.Send’den gelen mesajları yazar
func (c *Client) writePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

// ✅ ServeWs: bağlantı kurar, read/write pump başlatır
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// 🔑 Token doğrulama
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

	// 🔑 WebSocket upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		Conn: conn,
		Send: make(chan Message, 256), // buffered → yavaş client sistemi tıkamaz
		User: claims.Username,
		ID:   claims.UserID,
	}

	hub.Register <- client

	// Okuma ve yazma pump’larını başlat
	go client.readPump(hub)
	go client.writePump()
}
