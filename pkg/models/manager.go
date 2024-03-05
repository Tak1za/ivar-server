package models

import (
	"encoding/json"
	"log"
)

type Manager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func (m *Manager) send(message []byte, ignore *Client) {
	for conn := range m.Clients {
		//Send messages not to the shielded connection
		if conn != ignore {
			conn.Send <- message
		}
	}
}

func (m *Manager) Start() {
	for {
		select {
		case conn := <-m.Register:
			m.Clients[conn] = true
		case conn := <-m.Unregister:
			if _, ok := m.Clients[conn]; ok {
				close(conn.Send)
				delete(m.Clients, conn)
			}
		case msg := <-m.Broadcast:
			var jsonMsg Message
			if err := json.Unmarshal(msg, &jsonMsg); err != nil {
				log.Println("error converting message to correct format: " + err.Error())
				return
			}
			if jsonMsg.Recipient != "" {
				var clientToSendTo Client
				for client := range m.Clients {
					if client.Id == jsonMsg.Recipient {
						clientToSendTo = *client
					}
				}
				clientToSendTo.Send <- msg
			} else {
				for conn := range m.Clients {
					select {
					case conn.Send <- msg:
					default:
						close(conn.Send)
						delete(m.Clients, conn)
					}
				}
			}
		}
	}
}
