package websockets

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func SendMap(data map[string]interface{}) error {

	if conn == nil {
		return fmt.Errorf("websocket conexion not stablished")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// se crea la seccion critica para escribir el mensaje en el websocket
	mu.Lock()
	defer mu.Unlock()
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return err
	}

	return nil
}
