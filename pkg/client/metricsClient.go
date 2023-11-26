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
	//por cada grupo de metricas recibido cada cierto tiempo se actualizara la base de datos para tener un registro de el servidor
	//hay que tener en cuenta que en call metrics se va a quedar un proceso un tiempo infinito recibiendo cosas del servidor
	//asique la actualizacion de la base de datos se haria ahi, o simplemente hacerlo todo junto en runnodemetrics y ya.
	select {
	case <-ctx.Done():
		log.Println("Client: Graceful shutdown requested. Exiting runNodeCheck...")
		return
	default:

		for {
			conn, err := grpc.Dial(nodeAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Printf("client could not connect to %s: %v", nodeAddress, err)
			}

			client := pb.NewMetricServiceClient(conn)

			metrics := callMetrics(client, id, nodeAddress)

			conn.Close()

			if metrics == nil {
				break
			}
		}
	}
}

func callMetrics(client pb.MetricServiceClient, id int, nodeAddress string) error {

	timeout := 10 * time.Second

	timer := time.NewTimer(timeout)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metrics := make(map[string]string)

	stream, err := client.RequestMetrics(ctx, &pb.MetricsRequest{Id: int64(id)})
	if err != nil {
		log.Printf("Error making RequestMetrics call: %v", err)
		return err
	}

	go func() {
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Client: End Of File detected, closing in timeout: %d err: %v", timeout, err)
				break
			}
			if err != nil {
				log.Printf("Client: could not Recv, error while streaming, closing in timeout: %d err: %v", timeout, err)
				break
			}
			for key, value := range data.Metrics {
				metrics[key] = value
			}
			timer = time.NewTimer(timeout)
			log.Printf("Client: received a data: %v", metrics)
			models.UpdateDatabaseMetrics(nodeAddress, metrics) //que las metricas del log se actualicen cada mucho mas tiempo para no saturar, guardarlas cada 5 segundos,
			// ahi ya calcular que si recibimos un mensaje cada 1/4 de segundo cuantos mensajes hay que no guardar y cual si
		}
	}()
	for {
		// Comprueba si el temporizador ha caducado
		select {
		case <-timer.C:
			log.Printf("Client: Timeout expired on node: %v", nodeAddress)
			return ctx.Err()
		default:
		}
	}

}
