package Server

import (
	model "Block-P/pkg/models"
	metrics "Block-P/pkg/server/metrics"
	pb "Block-P/proto" // pakages generated with .proto
	"log"
	"net"

	"google.golang.org/grpc"
)

func Server() error {

	lis, err := net.Listen(model.GlobalConfig.Protocol, model.GlobalConfig.PortAddress)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMetricServiceServer(grpcServer, metrics.InitMetricsServer())

	log.Printf("server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start the server: %v", err)
		return err
	}
	return nil

}
