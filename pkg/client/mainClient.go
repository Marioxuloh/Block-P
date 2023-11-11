package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	// pakages generated with .proto
)

type Config struct {
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	MaxConnections int    `json:"maxConnections"`
	DebugMode      bool   `json:"debugMode"`
	Id             int    `json:"id"`
}

var config Config

const (
	callInterval = 10 * time.Second
)

var (
	shutdownRequested = false
)

func main() {
	// Configurar el manejador de señales para manejar Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM) //señales SIGINT(os.Interrupt) y SIGTERM(syscall.SIGTERM)

	// Read config.json
	configPath := filepath.Join("..", "..", "config", "config.json")
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

	// Lista de direcciones de nodos
	nodeAddresses := []string{"localhost:8080", "localhost:8080",
		"localhost:8080", "localhost:8080",
		"localhost:8080", "localhost:8080",
		"localhost:8080", "localhost:8080",
		"localhost:8080", "localhost:8080"}

	// Configurar contexto para el cierre ordenado
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Esperar la señal de interrupción
		<-sigCh
		log.Println("A interrupt signal was received. Shutting down gracefully...")
		shutdownRequested = true
		cancel()
	}()

	for !shutdownRequested {

		var wg sync.WaitGroup

		for _, addr := range nodeAddresses {
			wg.Add(1)
			go func(addr string) {
				defer wg.Done()
				runNodeCheck(ctx, addr)
			}(addr)
		}

		time.Sleep(callInterval)
		wg.Wait()

	}
}
