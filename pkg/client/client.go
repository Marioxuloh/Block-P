package Client

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Node
type Node struct {
	Name string
	Addr string
}

var (
	shutdownRequested = false
)

func Client(nodes []Node, id int) {
	// Configurar el manejador de señales para manejar Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM) //señales SIGINT(os.Interrupt) y SIGTERM(syscall.SIGTERM)

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

	for _, Node := range nodes {
		wg.Add(1)
		go func(addr string, name string) {
			defer wg.Done()
			runNodeMetrics(ctx, addr, name, id)
		}(Node.Addr, Node.Name)
	}

	wg.Wait()

}
