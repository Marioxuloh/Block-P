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

	//aqui para aprovechar el hecho de haber entablado un stream, se podria hacer aqui la espera
	// de intervalo ya que el servidor va a esperar indefinidamente informacion y asi con subrutinas
	// podemos estar mandando informacion de metricas cada medio segundo al cliente
	//nos ahorramos que cada medio segundo se tenga que establecer otra vez un stream
	//si no no tendria mucho sentido utilizar streams xdd

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
	//aqui mandar directamente el map de metrics
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
