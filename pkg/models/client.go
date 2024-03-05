package models

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type IClient interface {
	Read(manager *Manager)
	Write()
}

type Client struct {
	Id     string
	Socket *websocket.Conn
	Send   chan []byte
}

func NewClient(id string, socket *websocket.Conn, send chan []byte) *Client {
	return &Client{
		Id:     id,
		Socket: socket,
		Send:   send,
	}
}

func (c *Client) Read(manager *Manager) {
	defer func() {
		manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}

		jsonMessage, _ := json.Marshal(&Message{Sender: c.Id, Content: string(message)})
		manager.Broadcast <- jsonMessage
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
