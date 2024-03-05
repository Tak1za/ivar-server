package main

import (
	"ivar/pkg/controller"
	"ivar/pkg/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	manager := models.Manager{
		Broadcast:  make(chan []byte),
		Register:   make(chan *models.Client),
		Unregister: make(chan *models.Client),
		Clients:    make(map[*models.Client]bool),
	}

	ctrl := controller.New(&manager)
	r.GET("/ws", ctrl.HandleConnections)

	go manager.Start()

	if err := r.Run(":8080"); err != nil {
		panic("error creating server: " + err.Error())
	}
}
