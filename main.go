package main

import (
	"context"
	"ivar/pkg/controller"
	"ivar/pkg/database"
	"ivar/pkg/models"
	"ivar/pkg/user"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	r := gin.Default()

	manager := models.Manager{
		Broadcast:  make(chan []byte),
		Register:   make(chan *models.Client),
		Unregister: make(chan *models.Client),
		Clients:    make(map[*models.Client]bool),
	}

	go manager.Start()

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("error connecting to database: " + err.Error())
	}
	defer conn.Close(context.Background())

	store := database.NewStore(conn)
	userService := &user.Service{Store: store}

	ctrl := controller.New(&manager, userService)
	r.GET("/ws", ctrl.HandleConnections)
	r.POST("/api/v1/users", ctrl.CreateUser)

	if err := r.Run(":8080"); err != nil {
		panic("error creating server: " + err.Error())
	}
}
