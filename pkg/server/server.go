package Server

import (
	"log"
	"net"

	pb "Block-P/proto" // pakages generated with .proto

	"google.golang.org/grpc"
)

var id int

type connectionServer struct {
	pb.ConnectionServiceServer
}

type metricsServer struct {
	pb.MetricServiceServer
}

func Server(protocol string, address string, nodeID int) {

	id = nodeID

	lis, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterConnectionServiceServer(grpcServer, &connectionServer{})
	//pb.RegisterMetricServiceServer(grpcServer, metricsServer{})
	log.Printf("server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

}
