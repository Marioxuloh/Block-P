package Client

import (
	models "Block-P/pkg/models"
	pb "Block-P/proto"
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runNodeMetrics(ctx context.Context, nodeAddress string, id int) {
	select {
	case <-ctx.Done():
		log.Println("Client: Graceful shutdown requested. Exiting runNodeCheck...")
		return
	default:
		conn, err := grpc.Dial(nodeAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("client could not connect to %s: %v", nodeAddress, err)
		}
		defer conn.Close()

		client := pb.NewMetricServiceClient(conn)

		metrics := callMetrics(client, id)

		models.UpdateDatabaseMetrics(nodeAddress, metrics) //aqui hacemos update de las metricas y mandamos la informacion al sockets del controller

	}
}

func callMetrics(client pb.MetricServiceClient, id int) (metrics map[string]string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	metrics = make(map[string]string)

	stream, err := client.RequestMetrics(ctx, &pb.MetricsRequest{Id: int64(id)})
	if err != nil {
		log.Printf("Error making RequestMetrics call: %v", err)
		return nil
	}

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("client could not Recv, error while streaming: %v", err)
			return nil
		}

		for key, value := range data.Metrics {
			metrics[key] = value
		}
	}
	log.Printf("Client: received a data: %v", metrics)
	//log.Printf("Client: Streaming finished")

	return metrics

}
