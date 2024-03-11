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
	SendFriendRequest(c *gin.Context)
	UpdateFriendRequest(c *gin.Context)
	GetFriendRequests(ctx *gin.Context)
	GetFriends(ctx *gin.Context)
	RemoveFriend(ctx *gin.Context)
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
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.Create(user.ID, user.Username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *controller) SendFriendRequest(ctx *gin.Context) {
	var friendRequest models.AddFriendRequest
	if err := ctx.BindJSON(&friendRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.AddFriend(&friendRequest); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error adding friend request"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *controller) UpdateFriendRequest(ctx *gin.Context) {
	var friendRequest models.UpdateFriendRequest
	if err := ctx.BindJSON(&friendRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.UpdateFriend(&friendRequest); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating friend request"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *controller) GetFriendRequests(ctx *gin.Context) {
	userId, _ := ctx.Params.Get("userId")

	friendRequests, err := c.userService.GetFriendRequests(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting friend requests"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": friendRequests})
}

func (c *controller) GetFriends(ctx *gin.Context) {
	userId, _ := ctx.Params.Get("userId")

	friends, err := c.userService.GetFriends(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting friends"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": friends})
}

func (c *controller) RemoveFriend(ctx *gin.Context) {
	var deleteFriend models.RemoveFriendRequest
	if err := ctx.BindJSON(&deleteFriend); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := c.userService.RemoveFriend(deleteFriend.CurrentUserId, deleteFriend.ToRemoveUserId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting friend"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *controller) HandleConnections(ctx *gin.Context) {
	withUser, _ := ctx.Params.Get("userId")
	if withUser == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no user id provided"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	client := models.NewClient(withUser, conn, make(chan []byte))
	c.manager.Register <- client

	go client.Read(c.manager)
	go client.Write()
}
