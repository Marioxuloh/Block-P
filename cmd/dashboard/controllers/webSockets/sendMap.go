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

	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return err
	}

	return nil
}
