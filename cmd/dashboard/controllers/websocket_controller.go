package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketController struct {
	upgrader *websocket.Upgrader
	// Puedes agregar más campos según sea necesario
}

// Nuevo constructor para WebsocketController
func NewWebsocketController() *WebsocketController {
	return &WebsocketController{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

var activeConnections []*websocket.Conn //guardar los websockets activos para enviar msjes

func (wc *WebsocketController) HandleWebSocket(c *gin.Context) {

	var messageType int

	conn, err := wc.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Lógica de manejo de mensajes WebSocket

	messageType, _, err = conn.ReadMessage()
	if err != nil {
		return
	}

	// Aquí puedes agregar la lógica para enviar datos al cliente
	// En este ejemplo, enviamos un mensaje

	err = conn.WriteMessage(messageType, []byte("Hola, desde el servidor"))
	if err != nil {
		return
	}

}

// Función para enviar mensajes al dashboard
func (wc *WebsocketController) SendMessage(c *gin.Context) {
	message := "¡Hola desde el controlador!"
	wc.sendMessageToDashboard(message)
	c.JSON(http.StatusOK, gin.H{"message": "Mensaje enviado al dashboard"})
}

// Función para enviar mensajes al dashboard
func (wc *WebsocketController) sendMessageToDashboard(message string) {
	// Puedes almacenar las conexiones websocket activas y enviar mensajes a todas ellas
	// En este ejemplo, se envía el mensaje a todas las conexiones activas, pero en una aplicación real, deberías gestionar las conexiones activas de manera adecuada.
	for _, conn := range activeConnections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			// Manejar errores de escritura si es necesario
			continue
		}
	}
}
