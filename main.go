package main //para hacer el ejecutable hay que utilizar filepath.Abs y entregarlo con
// las dependencias las cuales son todo codigo ajeno a .go, como la plantilla html, el archivo config.json
// hay que cambiar en todos los sitios donde obtenemos estas dependencias de forma normal a ruta absoluta (de forma dinamica obviamente xddd)

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	dashboard "Block-P/cmd/dashboard"
	client "Block-P/pkg/client"
	server "Block-P/pkg/server"
)

//config

type Config struct {
	Port           int           `json:"port"`
	DashPort       int           `json:"dashPort"`
	Protocol       string        `json:"protocol"`
	MaxConnections int           `json:"maxConnections"`
	DebugMode      bool          `json:"debugMode"`
	Id             int           `json:"id"`
	MasterMode     bool          `json:"masterMode"` //tener en cuenta si es master o no que demomento no lo has tenido en cuenta en el codigo
	CallInterval   time.Duration `json:"callInterval"`
}

//server

var address string
var dashAddress string

var config Config

func main() {

	// Manejar señales de interrupción (SIGINT o CTRL+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{}) //canal por el que se notifica que han acabado las subrutinas

	//cierre ordenado
	var wg sync.WaitGroup

	// Read config.json
	configPath := filepath.Join("config.json")
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	address = ":" + strconv.Itoa(config.Port)
	dashAddress = ":" + strconv.Itoa(config.DashPort)

	//servidor

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Server(config.Protocol, address, config.Id)
	}()

	//cliente

	config.CallInterval = config.CallInterval * time.Second

	wg.Add(1)
	go func() {
		defer wg.Done()
		client.Client(config.CallInterval, config.Id)
	}()

	//dashboard

	wg.Add(1)
	go func() {
		defer wg.Done()
		dashboard.Dashboard(dashAddress)
	}()

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
