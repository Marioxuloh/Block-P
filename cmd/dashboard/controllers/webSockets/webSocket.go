package websockets

import (
	model "Block-P/pkg/models"
	"log"
	"net/http"
)

func WebSocketInit() error {

	http.HandleFunc("/ws", handleWebSocket)

	log.Printf("Initializing websocket server on :%s", model.GlobalConfig.WebSocketAddress)
	err := http.ListenAndServe(model.GlobalConfig.WebSocketAddress, nil)
	if err != nil {
		log.Printf("Error listen and serve websocket :", err)
		return err
	}
	return nil
}
