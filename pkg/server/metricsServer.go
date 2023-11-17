package Server

import (
	pb "Block-P/proto" // pakages generated with .proto
	"log"
)

func (s *metricsServer) RequestMetrics(req *pb.MetricServiceServer, stream pb.MetricService_RequestMetricsServer) error {
	log.Printf("Server: MetricsRequest received from nodeID: %v", req)

	response := &pb.Data{
		Id:      int64(id),
		Metrics: map[string]string{"cpu": "75"},
	}

	if err := stream.Send(response); err != nil {
		log.Printf("Error sending response: %v", err)
		return err
	}

	return nil

}
