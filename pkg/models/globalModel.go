package Model

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Node
type Node struct {
	Name string
	Addr string
}

// config
type Config struct {
	Port           int
	DashPort       int
	Protocol       string
	MaxConnections int
	DebugMode      bool
	ID             int64
	MasterMode     bool
	Secure         bool
	Name           string
	Ip             string
	PortAddress    string
	DashAddress    string
	FullAddress    string
	Nodes          []Node
}

var GlobalConfig Config

func InitGlobalData() error {
	// Obtener el directorio de configuración del usuario
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error al buscar archivos de configuracion:", err)
		return err
	}

	// Crear la ruta completa del directorio de configuración
	configDirPath := filepath.Join(homeDir, ".config", "block-p")

	// Verificar si el directorio de configuración existe
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		// Si no existe, crear el directorio
		err := os.MkdirAll(configDirPath, 0755)
		if err != nil {
			fmt.Println("Error al crear el directorio de configuración:", err)
			return err
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
			return err
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
name=master
ip=localhost

[nodes]

master=localhost:8080
`

		// Escribir los valores por defecto en el archivo
		_, err = file.WriteString(defaultConfig)
		if err != nil {
			fmt.Println("Error al escribir en el archivo de configuración:", err)
			file.Close()
			return err
		}

		file.Close()

		fmt.Printf("Se ha creado el archivo de configuración en: %v\n", configFilePath)
	}

	// Abrir el archivo config.config
	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return err
	}
	defer file.Close()

	// Crear un lector de líneas para el archivo
	scanner := bufio.NewScanner(file)

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
					fmt.Sscanf(value, "%d", &GlobalConfig.Port)
				case "dashPort":
					fmt.Sscanf(value, "%d", &GlobalConfig.DashPort)
				case "protocol":
					GlobalConfig.Protocol = value
				case "maxConnections":
					fmt.Sscanf(value, "%d", &GlobalConfig.MaxConnections)
				case "debugMode":
					fmt.Sscanf(value, "%t", &GlobalConfig.DebugMode)
				case "id":
					fmt.Sscanf(value, "%d", &GlobalConfig.ID)
				case "masterMode":
					fmt.Sscanf(value, "%t", &GlobalConfig.MasterMode)
				case "secure":
					fmt.Sscanf(value, "%t", &GlobalConfig.Secure)
				case "name":
					GlobalConfig.Name = value
				case "ip":
					GlobalConfig.Ip = value
				}
			case "nodes":
				switch key {
				default:
					// Asignar el valor correspondiente a la estructura Node
					nodeName := key
					nodeAddr := value
					GlobalConfig.Nodes = append(GlobalConfig.Nodes, Node{Name: nodeName, Addr: nodeAddr})
				}
			}
		}
	}

	// Verificar errores del escaneo
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return err
	}

	GlobalConfig.PortAddress = ":" + strconv.Itoa(GlobalConfig.Port)
	GlobalConfig.DashAddress = ":" + strconv.Itoa(GlobalConfig.DashPort)
	GlobalConfig.FullAddress = GlobalConfig.Ip + ":" + strconv.Itoa(GlobalConfig.Port)

	return nil
}
