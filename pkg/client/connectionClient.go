package Client

import (
	"context"
	"log"
	"time"

	models "Block-P/pkg/models"
	pb "Block-P/proto" // pakages generated with .proto

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runNodeCheck(ctx context.Context, nodeAddress string, id int) {

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

		client := pb.NewConnectionServiceClient(conn)

		connected := callConnection(client, id)

		models.UpdateDatabaseConnected(nodeAddress, connected)

	}

}

func callConnection(client pb.ConnectionServiceClient, id int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.RequestConnection(ctx, &pb.ConnectionRequest{Id: int64(id)})
	if err != nil {
		return false
	}
	log.Printf("Client: Acknowledge response from Id: %v message: %s", res.Id, res.Message)

	return true

}
