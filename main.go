package main

import (
	"context"
	"encoding/json"
	"fmt"
	"ivar/pkg/chat"
	"ivar/pkg/controller"
	"ivar/pkg/database"
	"ivar/pkg/models"
	"ivar/pkg/user"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
	Clients     map[*Client]bool
	Broadcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
	ChatService *chat.Service
}

var manager = Manager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
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
			var jsonMsg models.Message
			if err := json.Unmarshal(msg, &jsonMsg); err != nil {
				log.Println("error converting message to correct format: " + err.Error())
				return
			}
			if jsonMsg.Recipient != "" {
				if err := m.ChatService.AddMessage(jsonMsg); err != nil {
					log.Println("error adding message: " + err.Error())
				}
				var clientToSendTo Client
				for client := range m.Clients {
					if client.Id == jsonMsg.Recipient {
						clientToSendTo = *client
					}
				}
				if clientToSendTo.Id != "" {
					clientToSendTo.Send <- msg
				}
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

func (c *Client) Read() {
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

		var jsonMsg models.Message
		if err := json.Unmarshal(message, &jsonMsg); err != nil {
			log.Println("error converting message to correct format: " + err.Error())
			return
		}

		manager.Broadcast <- message
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnections(ctx *gin.Context) {
	currentUser, _ := ctx.Params.Get("userId")
	if currentUser == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no user id provided"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	client := &Client{
		Id:     currentUser,
		Socket: conn,
		Send:   make(chan []byte),
	}

	manager.Register <- client

	go client.Read()
	go client.Write()
}

func main() {
	r := gin.New()
	allowedOriginsFromEnv := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := strings.Split(allowedOriginsFromEnv, ",")
	r.Use(cors.New(cors.Config{
		AllowWildcard: true,
		AllowOrigins:  allowedOrigins,
		AllowMethods:  []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:  []string{"*"},
	}))

	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("error connecting to database: " + err.Error())
	}
	defer conn.Close()

	store := database.NewStore(conn)
	userService := &user.Service{Store: store}
	chatService := &chat.Service{Store: store}
	manager.ChatService = chatService

	go manager.Start()

	ctrl := controller.New(userService, chatService)
	r.GET("/ws/:userId", HandleConnections)
	r.POST("/api/v1/users", ctrl.CreateUser)
	r.POST("/api/v1/friends", ctrl.SendFriendRequest)
	r.PUT("/api/v1/friends", ctrl.UpdateFriendRequest)
	r.GET("/api/v1/friends/requests/:userId", ctrl.GetFriendRequests)
	r.GET("/api/v1/friends/:userId", ctrl.GetFriends)
	r.DELETE("/api/v1/friends", ctrl.RemoveFriend)
	r.POST("/api/v1/chats/info", ctrl.GetChatInfo)
	r.GET("/api/v1/chats/:userId", ctrl.GetAllChats)

	if err := r.Run(":8080"); err != nil {
		panic("error creating server: " + err.Error())
	}
}
