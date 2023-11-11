package main

import (
	"context"
	"log"
	"time"

	pb "Block-P/proto" // pakages generated with .proto

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

	//actualizar estado del nodo en la base de datos para que se vea reflejado en el dashboard

}

func callConnection(client pb.ConnectionServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.RequestConnection(ctx, &pb.ConnectionRequest{Id: int64(config.Id)})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	log.Printf("Acknowledge received from nodeID: %v message: %s", res.Id, res.Message)
}
