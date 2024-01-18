package websockets

import (
	model "Block-P/pkg/models"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func WebSocketInit() error {

	http.HandleFunc("/ws", handleWebSocket)

	// Configura CORS con opciones detalladas
	headers := handlers.AllowedHeaders([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:8081"}) // Reemplaza con tu URL de React
	handler := handlers.CORS(headers, methods, origins)(http.DefaultServeMux)

	log.Printf("Initializing websocket server on :%s", model.GlobalConfig.WebSocketAddress)
	err := http.ListenAndServe(model.GlobalConfig.WebSocketAddress, handler)
	if err != nil {
		log.Printf("Error listen and serve websocket: %v", err)
		return err
	}
	return nil
}
