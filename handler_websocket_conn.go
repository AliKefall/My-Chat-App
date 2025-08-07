package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"Message"`
}

var (
	clients = make(map[*websocket.Conn]bool)

	broadcast = make(chan Message)

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func (cfg *apiConfig) handlerWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrader Error!: ", err)
		return
	}

	defer func() {
		conn.Close()
		delete(clients, conn)
	}()

	clients[conn] = true

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Reading Error: ", err)
			break
		}
		broadcast <- msg
	}

}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println("Writing Error: ", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
