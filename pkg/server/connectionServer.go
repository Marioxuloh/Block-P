package Server

import (
	pb "Block-P/proto" // pakages generated with .proto
	"context"
	"log"
)

func (c *connectionServer) RequestConnection(ctx context.Context, req *pb.ConnectionRequest) (*pb.Acknowledge, error) {
	log.Printf("Server: ConnectionRequest received from nodeID: %v", req.Id)

	return &pb.Acknowledge{
		Id:      int64(id),
		Message: "ack",
	}, nil
}
