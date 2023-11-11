package main

import (
	"context"
	"log"
	"time"

	pb "Block-P/proto" // pakages generated with .proto
)

func callConnection(client pb.ConnectionServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.RequestConnection(ctx, &pb.ConnectionRequest{}) //agregar el id del cliente para que el servidor que reciba sepa el id
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	log.Printf("%s", res.Message)
}
