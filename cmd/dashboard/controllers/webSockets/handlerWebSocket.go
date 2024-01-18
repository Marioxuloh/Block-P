package websockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var conn *websocket.Conn

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	var err error
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error updating conexion:", err)
		return
	}
	defer conn.Close()

	log.Printf("Client websocket connected")

	for {
		// Maneja los mensajes del cliente si es necesario
		// se mantiene eternamente hasta que el cliente cierra la conexion
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Client websocket disconnected")
			break
		}
	}
}
