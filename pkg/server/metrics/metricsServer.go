package MetricsServer

import (
	pb "Block-P/proto" // pakages generated with .proto
	"log"
)

var Id int

type metricsServer struct {
	pb.MetricServiceServer
}

func InitMetricsServer() *metricsServer {
	return &metricsServer{}
}

func (s *metricsServer) RequestMetrics(req *pb.MetricsRequest, stream pb.MetricService_RequestMetricsServer) error {
	log.Printf("Server: MetricsRequest received from nodeID: %v", req.Id)

	cpu, err := getCPU()
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		cpu = "N/A"
	}
	mem, err := getMEM()
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		mem = "N/A"
	}

	metrics := map[string]string{
		"cpu":           cpu,
		"mem":           mem,
		"ram":           "60",
		"disk":          "80",
		"network":       "1000",
		"response_time": "50",
	}

	for metricName, metricValue := range metrics {
		response := &pb.Data{
			Id:      int64(Id),
			Metrics: map[string]string{metricName: metricValue},
		}

		if err := stream.Send(response); err != nil {
			log.Printf("Error sending response for metric %s: %v", metricName, err)
			return err
		}
	}

	return nil
}