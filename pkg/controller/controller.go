package controller

import (
	"fmt"
	"ivar/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Controller interface {
	HandleConnections(c *gin.Context)
}

type controller struct {
	manager *models.Manager
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func New(manager *models.Manager) *controller {
	return &controller{
		manager: manager,
	}
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
