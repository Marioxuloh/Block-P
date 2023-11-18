package Client

import (
	models "Block-P/pkg/models"
	pb "Block-P/proto"
	"context"
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

		models.UpdateDatabaseMetrics(nodeAddress, metrics) //podria cambiarlo y pasarle al mainclient un map querelacione nodo
		//y estado para hacer 1 sola escritura de todos los estados a la vez

	}
}

func callMetrics(client pb.MetricServiceClient, id int) (metrics map[string]string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := client.RequestMetrics(ctx, &pb.MetricsRequest{Id: int64(id)})
	if err != nil {
		return nil
	}
	log.Printf("Client: Data response from Id: %v", res)

	return metrics

}
