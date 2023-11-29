package main //para hacer el ejecutable hay que utilizar filepath.Abs y entregarlo con
// las dependencias las cuales son todo codigo ajeno a .go, como la plantilla html, el archivo config.json
// hay que cambiar en todos los sitios donde obtenemos estas dependencias de forma normal a ruta absoluta (de forma dinamica obviamente xddd)

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	dashboard "Block-P/cmd/dashboard"
	client "Block-P/pkg/client"
	server "Block-P/pkg/server"
)

// config
type Config struct {
	Port           int
	DashPort       int
	Protocol       string
	MaxConnections int
	DebugMode      bool
	ID             int
	MasterMode     bool
	Secure         bool
	Nodes          []Node
}

// Node estructura para almacenar información de un nodo
type Node struct {
	Name string
	Addr string
}

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

	// Abrir el archivo config.config
	file, err := os.Open("config.config")
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Crear un lector de líneas para el archivo
	scanner := bufio.NewScanner(file)

	// Crear una instancia de Config
	var config Config

	// Variable para determinar la sección actual del archivo
	var currentSection string

	// Leer el archivo línea por línea
	for scanner.Scan() {
		line := scanner.Text()
		// Verificar si la línea representa una sección
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
			continue
		}
		// Dividir la línea en clave y valor
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Asignar el valor correspondiente a la estructura Config
			switch currentSection {
			case "config":
				switch key {
				case "port":
					fmt.Sscanf(value, "%d", &config.Port)
				case "dashPort":
					fmt.Sscanf(value, "%d", &config.DashPort)
				case "protocol":
					config.Protocol = value
				case "maxConnections":
					fmt.Sscanf(value, "%d", &config.MaxConnections)
				case "debugMode":
					fmt.Sscanf(value, "%t", &config.DebugMode)
				case "id":
					fmt.Sscanf(value, "%d", &config.ID)
				case "masterMode":
					fmt.Sscanf(value, "%t", &config.MasterMode)
				case "secure":
					fmt.Sscanf(value, "%t", &config.Secure)
				}
			case "nodes":
				switch key {
				default:
					// Asignar el valor correspondiente a la estructura Node
					nodeName := key
					nodeAddr := value
					config.Nodes = append(config.Nodes, Node{Name: nodeName, Addr: nodeAddr})
				}
			}
		}
	}

	// Verificar errores del escaneo
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	address = ":" + strconv.Itoa(config.Port)
	dashAddress = ":" + strconv.Itoa(config.DashPort)

	//servidor, gestionamos todas las llamadas entrantes

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Server(config.Protocol, address, config.ID)
	}()

	//cliente, gestionamos todos los mensajes que vamos a enviar: requestMetrics()

	wg.Add(1)
	go func() {
		defer wg.Done()
		client.Client(config.ID)
	}()

	//dashboard, desplegamos unh dashboard para visualizar los nodos, su informacion y poder inyectar codigo en ellos para utilizarlos como microservicios

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
