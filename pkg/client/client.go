package Client

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	MaxConnections int    `json:"maxConnections"`
	DebugMode      bool   `json:"debugMode"`
	Id             int    `json:"id"`
}

//var config Config

var (
	shutdownRequested = false
)

func Client(id int) {
	// Configurar el manejador de señales para manejar Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM) //señales SIGINT(os.Interrupt) y SIGTERM(syscall.SIGTERM)

	// Lista de direcciones de nodos
	nodeAddresses := []string{"localhost:8080"}

	// Configurar contexto para el cierre ordenado
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Esperar la señal de interrupción
		<-sigCh
		log.Println("Client: A interrupt signal was received. Shutting down gracefully...")
		shutdownRequested = true
		cancel()
	}()

	var wg sync.WaitGroup

	for _, addr := range nodeAddresses {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			//runNodeCheck(ctx, addr, id)
			runNodeMetrics(ctx, addr, id)
		}(addr)
	}

	wg.Wait()

}
