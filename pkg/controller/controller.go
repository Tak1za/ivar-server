package controller

import (
	"fmt"
	"ivar/pkg/models"
	"ivar/pkg/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Controller interface {
	HandleConnections(c *gin.Context)
	CreateUser(c *gin.Context)
}

type controller struct {
	manager     *models.Manager
	userService *user.Service
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func New(manager *models.Manager, userService *user.Service) *controller {
	return &controller{
		manager:     manager,
		userService: userService,
	}
}

func (c *controller) CreateUser(ctx *gin.Context) {
	var user models.User
	ctx.BindJSON(&user)

	if err := c.userService.Create(user.ID, user.Username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *controller) HandleConnections(ctx *gin.Context) {
	userId, _ := ctx.Params.Get("userId")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no user id provided"})
		return
	}

	for client := range c.manager.Clients {
		if client.Id == userId {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "connection already exists"})
			return
		}
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	client := models.NewClient(userId, conn, make(chan []byte))
	c.manager.Register <- client

	go client.Read(c.manager)
	go client.Write()
}
