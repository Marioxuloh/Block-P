package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	pb "Block-P/proto" // pakages generated with .proto

	"google.golang.org/grpc"
)

type Config struct {
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	MaxConnections int    `json:"maxConnections"`
	DebugMode      bool   `json:"debugMode"`
	Id             int    `json:"id"`
}

var config Config

type connectionServer struct {
	pb.ConnectionServiceServer
}

func main() {

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

	address := ":" + strconv.Itoa(config.Port)

	lis, err := net.Listen(config.Protocol, address)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterConnectionServiceServer(grpcServer, &connectionServer{})
	log.Printf("server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start: %v", err)
	}

}
