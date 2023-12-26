package MetricsClient

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

func MetricsRequestFromNodeToMaster(fullMasterAddress string, fullNodeAddress string, name string, id int64) error {

	conn, err := grpc.Dial(fullMasterAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Client: could not create dial to %v", fullMasterAddress)
		return err
	}
	defer conn.Close()

	client := pb.NewMetricServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := client.RequestMetricsFromNode(ctx, &pb.MetricsRequestTrigger{NodeAddress: string(fullNodeAddress), Name: string(name), Id: int64(id)})
	if err != nil {
		log.Printf("Client: could not send RequestMetricsFromNode to %v", fullMasterAddress)
		return err
	}

	log.Printf("Client: on RequestMetricsFromNode service, received %v from master", res)

	return nil
}

func RunNodeMetrics(nodeAddress string, name string, id int64, maxRetries int, timeout time.Duration, eachMetrics int) error {

	retries := 0
	for {
		conn, err := grpc.Dial(nodeAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Client: could not create dial to %v with name: %v: %v", nodeAddress, name, err)
			return err
		}

		client := pb.NewMetricServiceClient(conn)

		metrics := callMetrics(client, id, nodeAddress, name, timeout, eachMetrics)

		conn.Close()

		retries++

		if metrics == io.EOF {
			log.Printf("Client: on RequestMetrics service, closing streaming, EOF received from node: %v", nodeAddress) //el nodo ha mandado un eof
			return nil
		} else if retries >= maxRetries {
			log.Printf("Client: on RequestMetrics service, closing streaming, max retries detected on node: %v with err: %v", nodeAddress, metrics) //se han terminado los retires
			return nil
		} else {
			log.Printf("Client: on RequestMetrics service, retrying connect to node: %v with err: %v", nodeAddress, metrics) //retrying connect
		}

	}

}

func callMetrics(client pb.MetricServiceClient, id int64, nodeAddress string, name string, timeout time.Duration, eachMetrics int) error {

	var response error

	timer := time.NewTimer(timeout)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metrics := make(map[string]string)

	stream, err := client.RequestMetrics(ctx, &pb.MetricsRequest{Id: int64(id)})
	if err != nil {
		response = err
	} else {
		go func() {
			aux := 0
			for {
				data, err := stream.Recv()
				if err == io.EOF {
					log.Printf("Client: on RequestMetrics service, End Of File detected, closing streaming in timeout: %d err: %v", timeout, err)
					response = io.EOF //salimos totalmente y no volvemos a intentar conectar
					break
				}
				if err != nil {
					log.Printf("Client: on RequestMetrics service, could not Recv, error while streaming, retrying connect in timeout: %d err: %v", timeout, err)
					response = err //salimos y volvemos a intentar conectar pasados 13s
					timer = time.NewTimer(timeout)
					break
				}
				for key, value := range data.Metrics {
					metrics[key] = value
				}
				timer = time.NewTimer(timeout)
				aux++
				if aux == eachMetrics {
					models.UpdateDatabaseMetrics(nodeAddress, name, metrics) //esto cada 5s como fibonacci, si cada 1/4s llega un metrics cada 15 metrics uno se guarda en el log
					aux = 0
				}
				models.UpdateDashboardMetrics(nodeAddress, name, metrics) //cada 1/4s
			}
		}()
	}
	for {
		select {
		case <-timer.C:
			return response
		default:
		}
	}
}
