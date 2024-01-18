package ModelMetrics

import (
	websocket "Block-P/cmd/dashboard/controllers/webSockets"
	model "Block-P/pkg/models"
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func UpdateDatabaseMetrics(nodeAddress string, name string, metrics map[string]string) {
	// Lógica para actualizar la base de datos de las métricas de un nodo, en este caso, el estado de conexión.
	// Puedes acceder a los elementos del mapa dentro de la función.
	//for key, value := range metrics {
	// Hacer algo con key y value, por ejemplo, insertar en la base de datos.
	// key es la clave del mapa, y value es el valor asociado.
	//}
}

func UpdateDashboardMetrics(nodeAddress string, name string, metrics map[string]string) {
	// Añade nodeAddress y name al mapa de métricas
	metrics["nodeAddress"] = nodeAddress
	metrics["name"] = name

	// Construye un mapa con la información completa
	data := map[string]interface{}{
		"nodeAddress": nodeAddress,
		"name":        name,
		"metrics":     metrics,
	}

	// Envia el mapa a través del websocket
	err := websocket.SendMap(data)
	if err != nil {
		log.Printf("Error sending map by websockets: ", err)
		// Puedes manejar el error según tus necesidades
	}
}

func GetAddons() (map[string]string, error) {
	//buscar todos los .bp y ejecutar los scripts, devolver map
	//obtener la lista de archivos *.bp
	archives, err := getArchivesBP(model.GlobalConfig.RouteAddons)

	if err != nil {
		log.Println("Error Model:", err)
		return nil, err
	}

	metricsAddons := make(map[string]string)

	// Procesar cada archivo
	for _, archive := range archives {
		// Abrir el archivo en modo lectura
		file, err := os.Open(archive)
		if err != nil {
			log.Println("Error al abrir el archivo:", err)
			continue // Continuar con el siguiente archivo en caso de error
		}
		defer file.Close()

		// Crear un lector de líneas para el archivo
		scanner := bufio.NewScanner(file)

		// Variable para determinar la sección actual del archivo
		var currentSection string

		var name string
		var routeScript string

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
				case "service":
					switch key {
					case "name":
						name = value
					case "route":
						routeScript = value
					}
				}
			}

		}
		//implementar ejecucion del script y gardado de informacion en el map metricsAddons
		output, err := executeScript(routeScript)
		if err != nil {
			log.Fatal("Error executing script:", err)
			continue //salta al siguiente .bp
		}

		metricsAddons[name] = output
	}
	return metricsAddons, nil
}

func getArchivesBP(carpeta string) ([]string, error) {
	var archivos []string

	// Función anónima para filtrar archivos con extensión .bp
	filtro := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".bp") {
			archivos = append(archivos, path)
		}
		return nil
	}

	// Recorre la carpeta y aplica el filtro a cada archivo
	err := filepath.Walk(carpeta, filtro)
	if err != nil {
		return nil, err
	}

	return archivos, nil
}
func executeScript(scriptPath string) (string, error) {
	cmd := exec.Command(model.GlobalConfig.Shell, scriptPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing script: %v\n", err)
		return "N/A", err
	}

	// Convierte la salida del script a una cadena
	// Eliminar caracteres de nueva línea al final de la cadena
	outputString := strings.TrimRight(string(output), "\n")

	return outputString, nil
}
