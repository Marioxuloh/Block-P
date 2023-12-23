package main //para hacer el ejecutable hay que utilizar filepath.Abs y entregarlo con
// las dependencias las cuales son todo codigo ajeno a .go, como la plantilla html, el archivo config.json
// hay que cambiar en todos los sitios donde obtenemos estas dependencias de forma normal a ruta absoluta (de forma dinamica obviamente xddd)

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	dashboard "Block-P/cmd/dashboard"
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

	//servidor, gestionamos todas las llamadas entrantes
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.Server()
		if err != nil {
			log.Printf("Main: Server error %v", err)
			return
		}
	}()

	//cliente, gestionamos todos los mensajes que vamos a enviar: requestMetrics(), demomento solo para master
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
