package websockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	conn *websocket.Conn
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocketHandler maneja las conexiones WebSocket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Configura el Upgrader para permitir la conexión desde el mismo origen
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// Actualiza la conexión HTTP a una conexión WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al actualizar la conexión: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("Cliente WebSocket conectado")

	// Simplemente espera hasta que se cierre la conexión
	select {}
}
