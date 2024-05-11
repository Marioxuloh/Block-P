package websockets

import (
	model "Block-P/pkg/models"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func WebSocketInit() error {

	http.HandleFunc("/ws", handleWebSocket)

	//origins := handlers.AllowedOrigins([]string{"http://" + model.GlobalConfig.DashAddress})
	origins := handlers.AllowedOrigins([]string{"*"})
	handler := handlers.CORS(origins)(http.DefaultServeMux)

	log.Printf("Initializing websocket server on :%s", model.GlobalConfig.WebSocketAddress)
	err := http.ListenAndServe(model.GlobalConfig.WebSocketAddress, handler)
	if err != nil {
		log.Printf("Error listen and serve websocket: %v", err)
		return err
	}
	return nil
}
