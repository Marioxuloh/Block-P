package Server

import (
	"log"
	"net"

	metrics "Block-P/pkg/server/metrics"
	pb "Block-P/proto" // pakages generated with .proto

	"google.golang.org/grpc"
)

var id int

func Server(protocol string, address string, nodeID int) {

	id = nodeID

	metrics.Id = id

	lis, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMetricServiceServer(grpcServer, metrics.InitMetricsServer())
	log.Printf("server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

}
