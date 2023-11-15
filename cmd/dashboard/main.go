package main

import (
	"Block-P/cmd/dashboard/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Static("/dashboard", "./view")

	// Configurar rutas
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Rutas de websockets
	wsController := controllers.NewWebsocketController()
	router.GET("/ws", wsController.HandleWebSocket)
	router.GET("/send-message", wsController.SendMessage)

	router.Run(":8080")
}
