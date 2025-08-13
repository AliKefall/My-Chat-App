package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	ID    string
	User  string
	Color string
	Conn  *websocket.Conn
	Pool  *Pool
}

type Message struct {
	Type      int    `json:"type"`
	Body      string `json:"body"`
	User      string `json:"user"`
	Color     string `json:"color"`
	TimeStamp string `json:"timeStamp"`
}

type MessageData struct {
	Message string
	Id      string
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()
}
