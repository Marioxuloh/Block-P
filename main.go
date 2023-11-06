package main

import (
    "encoding/json"
    "log"
    "net"
    "os"

    "google.golang.org/grpc"
	pb "Block-P/api" // pakages generated with .proto
    "Block-P/pkg/server" // own pkg for service implementation for the entrance petitions
	"Block-P/pkg/client" // own pkg for service implementation for the exit petitions
)

type Config struct {
    Port          int    `json:"port"`
    Protocol      string `json:"protocol"`
    MaxConnections int    `json:"maxConnections"`
    DebugMode     bool   `json:"debugMode"`
}

func main() {

	// Read config.json
	configFile, err := os.Open("config/config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	// Define server gRPC
	server := grpc.NewServer()

	// Regist services gRPC defined in .proto
	pb.RegisterMiAplicacionServer(server, &server.MyServer{})


	address := ":" + strconv.Itoa(config.Port)

	// Address and port
	listener, err := net.Listen(config.Protocol, address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server is listening on %s", address)

	// start server
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}


}