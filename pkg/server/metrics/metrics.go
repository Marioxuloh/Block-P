package MetricsServer

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func getCPU() (string, error) {
	// Ejecutar el comando ps y capturar la salida
	numCores := runtime.NumCPU()
	cmd := exec.Command("ps", "-eo", "%cpu")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Procesar la salida para obtener el uso total de la CPU
	lines := strings.Split(string(output), "\n")
	totalCPU := 0.0

	for _, line := range lines[1:] { // Comenzar desde la segunda línea para omitir la cabecera
		if len(line) > 0 {
			cpuUsage, err := parser(line)
			if err != nil {
				return "", err
			}
			totalCPU += cpuUsage
		}
	}

	totalCPU = totalCPU / float64(numCores)

	return fmt.Sprintf("%.2f", totalCPU), nil
}

func getMEM() (string, error) {
	// Ejecutar el comando ps y capturar la salida
	cmd := exec.Command("ps", "-eo", "%mem")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Procesar la salida para obtener el uso total de la ram
	lines := strings.Split(string(output), "\n")
	totalMEM := 0.0

	for _, line := range lines[1:] { // Comenzar desde la segunda línea para omitir la cabecera
		if len(line) > 0 {
			memUsage, err := parser(line)
			if err != nil {
				return "", err
			}
			totalMEM += memUsage
		}
	}

	return fmt.Sprintf("%.2f", totalMEM), nil
}

func getDISK() (string, error) {
	// Ejecutar el comando "df" y capturar la salida
	cmd := exec.Command("df", "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Procesar la salida para obtener el porcentaje de uso del espacio en disco
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("output format not as expected")
	}

	// Inicializar variables para el total de porcentajes y el total de tamaños de sistemas de archivos
	var totalUsed float64
	var totalAvail float64

	// Iterar sobre las líneas de salida
	for _, line := range lines[1:] { // Comenzar desde la segunda línea para omitir la cabecera
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			// Extraer el tamaño utilizado y disponible del sistema de archivos
			usedStr := fields[2]
			availStr := fields[3]

			used, err := parseSize(usedStr)
			if err != nil {
				return "", err
			}

			avail, err := parseSize(availStr)
			if err != nil {
				return "", err
			}

			// Sumar el tamaño utilizado y disponible al total
			totalUsed += used
			totalAvail += avail
		}
	}

	// Calcular el porcentaje total del uso del espacio en disco
	if totalUsed+totalAvail > 0 {
		percentage := (totalUsed * 100) / (totalUsed + totalAvail)
		return fmt.Sprintf("%.2f", percentage), nil
	}

	return "0", nil
}

func parseSize(sizeStr string) (float64, error) {

	// Convierte tamaños de disco de formato ("64G") a megabytes
	sizeStr = strings.ReplaceAll(sizeStr, ",", ".")
	unit := sizeStr[len(sizeStr)-1]
	valueStr := sizeStr[:len(sizeStr)-2]

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, err
	}

	switch unit {
	case 'G':
		value = value * 1024 // Convertir gigabytes a megabytes
	default:

	}

	return value, nil
}

func parser(line string) (float64, error) {
	var n float64
	// Convertir la cadena a un número decimal, float64
	_, err := fmt.Sscanf(line, "%f", &n)
	if err != nil {
		return 0.0, err
	}

	return n, nil
}
