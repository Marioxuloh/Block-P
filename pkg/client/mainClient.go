package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "Block-P/proto" // pakages generated with .proto

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

func runNodeCheck(ctx context.Context, nodeAddress string) {

	select {
	case <-ctx.Done():
		log.Println("Graceful shutdown requested. Exiting runNodeCheck...")
		return
	default:
		conn, err := grpc.Dial(nodeAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("could not connect to %s: %v", nodeAddress, err)
			return
		}
		defer conn.Close()

		client := pb.NewConnectionServiceClient(conn)

		callConnection(client)
	}

}
