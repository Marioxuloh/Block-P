package main //para hacer el ejecutable hay que utilizar filepath.Abs y entregarlo con
// las dependencias las cuales son todo codigo ajeno a .go, como la plantilla html, el archivo config.json
// hay que cambiar en todos los sitios donde obtenemos estas dependencias de forma normal a ruta absoluta (de forma dinamica obviamente xddd)

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
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
	Nodes          []client.Node
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

	// Obtener el directorio de configuración del usuario
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error al buscar archivos de configuracion:", err)
		return
	}

	// Crear la ruta completa del directorio de configuración
	configDirPath := filepath.Join(homeDir, ".config", "block-p")

	// Verificar si el directorio de configuración existe
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		// Si no existe, crear el directorio
		err := os.MkdirAll(configDirPath, 0755)
		if err != nil {
			fmt.Println("Error al crear el directorio de configuración:", err)
			return
		}

		fmt.Printf("Se ha creado el directorio de configuración en: %v\n", configDirPath)
	}

	// Crear la ruta completa del archivo de configuración
	configFilePath := filepath.Join(configDirPath, "config.config")

	// Verificar si el archivo de configuración existe
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Si no existe, crear el archivo con valores por defecto
		file, err := os.Create(configFilePath)
		if err != nil {
			fmt.Println("Error al crear el archivo de configuración:", err)
			return
		}

		// Definir los valores por defecto
		defaultConfig := `[config]

port=8080
dashPort=8081
protocol=tcp
maxConnections=100
debugMode=true
id=0
masterMode=true
secure=false

[nodes]

master=localhost:8080
`

		// Escribir los valores por defecto en el archivo
		_, err = file.WriteString(defaultConfig)
		if err != nil {
			fmt.Println("Error al escribir en el archivo de configuración:", err)
			file.Close()
			return
		}

		file.Close()

		fmt.Printf("Se ha creado el archivo de configuración en: %v\n", configFilePath)
	}

	// Abrir el archivo config.config
	file, err := os.Open(configFilePath)
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
					config.Nodes = append(config.Nodes, client.Node{Name: nodeName, Addr: nodeAddr})
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

	if config.MasterMode == true {
		//cliente, gestionamos todos los mensajes que vamos a enviar: requestMetrics(), demomento solo para master
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.Client(config.Nodes, config.ID)
		}()

		//dashboard, desplegamos unh dashboard para visualizar los nodos, su informacion y poder inyectar codigo en ellos para utilizarlos como microservicios, solo para master
		wg.Add(1)
		go func() {
			defer wg.Done()
			dashboard.Dashboard(dashAddress)
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
