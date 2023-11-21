package MetricsServer

import (
	"fmt"
	"os/exec"
	"runtime"
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

func parser(line string) (float64, error) {
	var n float64
	// Convertir la cadena de uso de CPU a un número decimal
	_, err := fmt.Sscanf(line, "%f", &n)
	if err != nil {
		return 0.0, err
	}

	return n, nil
}
