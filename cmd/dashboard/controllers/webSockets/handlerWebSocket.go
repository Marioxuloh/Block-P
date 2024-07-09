package websockets

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	mu   sync.Mutex
	conn *websocket.Conn
)

// WebSocketHandler maneja las conexiones WebSocket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Configura el Upgrader para permitir la conexión desde el mismo origen
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// Actualiza la conexión HTTP a una conexión WebSocket
	var err error
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al actualizar la conexión: %v", err)
		return
	}
	defer conn.Close()

	// Simplemente espera hasta que se cierre la conexión
	select {}
}
