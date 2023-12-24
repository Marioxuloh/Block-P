package MetricsServer

import (
	metricsClient "Block-P/pkg/client/metrics"
	model "Block-P/pkg/models"
	pb "Block-P/proto" // pakages generated with .proto
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type metricsServer struct {
	pb.MetricServiceServer
}

func InitMetricsServer() *metricsServer {
	return &metricsServer{}
}

func (s *metricsServer) RequestMetrics(req *pb.MetricsRequest, stream pb.MetricService_RequestMetricsServer) error {

	log.Printf("Server: MetricsRequest received from nodeID: %v", req.Id)

	// Create a channel to listen for interrupts
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{}) //canal por el que se notifica que han acabado las subrutinas

	// Goroutine to handle interrupts and close the stream
	go func() {
		<-interrupt
		log.Println("Server: on RequestMetrics service, A interrupt signal was received. Closing the stream.")
		// Close the stream gracefully by sending an EOF
		if err := stream.Send(nil); err != nil {
			log.Printf("Server: on RequestMetrics service, Error sending EOF: %v", err)
		}
		close(done)
	}()

	for {

		select {
		case <-done:
			return nil
		default:
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
				Id:      model.GlobalConfig.ID,
				Metrics: metrics,
			}

			if err := stream.Send(response); err != nil {
				log.Printf("Server: on RequestMetrics service, Error sending response for metric %s: %v", metrics, err)
			}

			time.Sleep(time.Second / 4) //se envian metricas cada 1/4 de segundo
		}
	}
}

func (s *metricsServer) RequestMetricsFromNode(ctx context.Context, req *pb.MetricsRequestTrigger) (*pb.Ack, error) {
	log.Printf("Server: on RequestMetricsFromNode service, received MetricsRequestTrigger from nodeAddress: %v name: %v nodeID: %v", req.NodeAddress, req.Name, req.Id)
	response := &pb.Ack{
		Ack: string("success"),
	}
	go metricsClient.RunNodeMetrics(req.NodeAddress, req.Name, req.Id, 5, 13*time.Second, 15)
	return response, nil
}
