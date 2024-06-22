package controller

import (
	"ivar/pkg/chat"
	"ivar/pkg/models"
	"ivar/pkg/server"
	"ivar/pkg/user"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateUser(c *gin.Context)
	SendFriendRequest(c *gin.Context)
	UpdateFriendRequest(c *gin.Context)
	GetFriendRequests(ctx *gin.Context)
	GetFriends(ctx *gin.Context)
	RemoveFriend(ctx *gin.Context)
	GetChatInfo(ctx *gin.Context)
	AddMessage(ctx *gin.Context)
	GetMessages(ctx *gin.Context)
	GetAllChats(ctx *gin.Context)
	CreateServer(ctx *gin.Context)
	GetServers(ctx *gin.Context)
	CreateInvite(ctx *gin.Context)
}

type controller struct {
	userService   *user.Service
	chatService   *chat.Service
	serverService *server.Service
}

func New(userService *user.Service, chatService *chat.Service, serverService *server.Service) *controller {
	return &controller{
		userService:   userService,
		chatService:   chatService,
		serverService: serverService,
	}
}

func (c *controller) CreateServer(ctx *gin.Context) {
	var server models.Server
	if err := ctx.BindJSON(&server); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.serverService.CreateServer(server.Name, server.OwnerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating server"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *controller) GetServers(ctx *gin.Context) {
	res, err := c.serverService.GetServers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting servers"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": res})
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

func (c *controller) GetChatInfo(ctx *gin.Context) {
	var chatInfoRequest models.ChatInfoRequest
	if err := ctx.BindJSON(&chatInfoRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	chatInfo, err := c.userService.GetChatInfo(chatInfoRequest.Users)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting chat info"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": chatInfo})
}

func (c *controller) AddMessage(ctx *gin.Context) {
	var message models.Message
	if err := ctx.BindJSON(&message); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := c.chatService.AddMessage(message); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *controller) GetAllChats(ctx *gin.Context) {
	userId, _ := ctx.Params.Get("userId")

	chats, err := c.chatService.GetAllChats(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting chats"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": chats})
}

func (c *controller) CreateInvite(ctx *gin.Context) {
	serverId, _ := ctx.Params.Get("serverId")

	id, err := strconv.Atoi(serverId)
	if err != nil {
		log.Println("bad serverId: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad server id"})
		return
	}

	inviteCode, err := c.serverService.CreateInvite(id)
	if err != nil {
		if err.Error() == "invalid server" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid server"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error generating invite"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": inviteCode})
}
