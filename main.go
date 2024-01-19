package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	dashboard "Block-P/cmd/dashboard"
	websocket "Block-P/cmd/dashboard/controllers/webSockets"
	client "Block-P/pkg/client"
	model "Block-P/pkg/models"
	server "Block-P/pkg/server"
)

func main() {
	// Manejar señales de interrupción (SIGINT o CTRL+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{}) //canal por el que se notifica que han acabado las subrutinas

	err := model.InitGlobalData()
	if err != nil {
		log.Printf("Main: init global data error %v", err)
		return
	}

	//cierre ordenado
	var wg sync.WaitGroup

	//servidor, gestionamos todas las llamadas entrantes, si es master espera a que se establezca la conexion websocket(llamamos al cliente y este hace que se envien cosas por websockets)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.Server()
		if err != nil {
			log.Printf("Main: Server error %v", err)
			return
		}
	}()

	//cliente, gestionamos todos los mensajes que vamos a enviar: requestMetrics(), si es master espera a que se establezca la conexion websocket(en el cliente enviamos mensajes)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := client.Client()
		if err != nil {
			log.Printf("Main: Client error %v", err)
			return
		}
	}()

	if model.GlobalConfig.MasterMode == true {
		//websocket, inicializamos el websocket para comunicarnos en tiempo real con el dashboard
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := websocket.WebSocketInit()
			if err != nil {
				log.Printf("Main: Websocket error %v", err)
				return
			}
		}()
		//dashboard, desplegamos unh dashboard para visualizar los nodos, su informacion y poder inyectar codigo en ellos para utilizarlos como microservicios, solo para master
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := dashboard.Dashboard()
			if err != nil {
				log.Printf("Main: Dashboard error %v", err)
				return
			}
		}()
	}
	//cierre ordenado
	go func() {
		// Esperar señales
		<-interrupt

		// Cerrar los canales y esperar que todas las goroutines finalicen
		close(done)
		wg.Wait()

		// Salir del programa
		os.Exit(0)
	}()
	<-done
	log.Println("Main: Received interrupt signal. Exiting...")

}
