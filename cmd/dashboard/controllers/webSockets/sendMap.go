package websockets

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// SendMap envía un mapa a través de la conexión websocket.
func SendMap(data map[string]interface{}) error {

	if conn == nil {
		return fmt.Errorf("websocket conexion not stablished")
	}

	// Serializa el mapa a formato JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Envía el JSON a través del websocket
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return err
	}

	return nil
}
