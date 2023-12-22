package MetricsServer

import (
	pb "Block-P/proto" // pakages generated with .proto
	"log"
	"time"
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

	for {
		cpu, err := getCPU()
		if err != nil {
			log.Printf("Error getting CPU usage: %v", err)
			cpu = "N/A"
		}
		mem, err := getMEM()
		if err != nil {
			log.Printf("Error getting MEM usage: %v", err)
			mem = "N/A"
		}
		disk, err := getDISK()
		if err != nil {
			log.Printf("Error getting DISK usage: %v", err)
			disk = "N/A"
		}

		metrics := map[string]string{
			"cpu":  cpu,
			"mem":  mem,
			"disk": disk,
		}

		response := &pb.Data{
			Id:      int64(Id),
			Metrics: metrics,
		}

		if err := stream.Send(response); err != nil {
			log.Printf("Server: Error sending response for metric %s: %v", metrics, err)
			return err
		}

		log.Printf("Server: send response in requestMetrics streaming: %v", response.Metrics)

		time.Sleep(time.Second / 4) //se envian metricas cada 1/4 de segundo
	}

}
