package main

import (
	pb "Block-P/proto" // pakages generated with .proto
	"context"
	"log"
)

func (c *connectionServer) RequestConnection(ctx context.Context, req *pb.ConnectionRequest) (*pb.Acknowledge, error) {
	log.Printf("RequestConnection received")
	return &pb.Acknowledge{
		Message: "ack",
	}, nil
}
