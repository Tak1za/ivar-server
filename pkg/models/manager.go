package models

import "encoding/json"

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
			jsonMessage, _ := json.Marshal(&Message{
				Content: "A new socket has connected.",
			})
			m.send(jsonMessage, conn)
		case conn := <-m.Unregister:
			if _, ok := m.Clients[conn]; ok {
				close(conn.Send)
				delete(m.Clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "A socket has disconnected."})
				m.send(jsonMessage, conn)
			}
		case msg := <-m.Broadcast:
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
