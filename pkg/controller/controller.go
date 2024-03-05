package controller

import (
	"fmt"
	"ivar/pkg/models"
	"ivar/pkg/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
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
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	client := models.NewClient(uuid.Must(uuid.NewV4(), nil).String(), conn, make(chan []byte))
	c.manager.Register <- client

	go client.Read(c.manager)
	go client.Write()
}
