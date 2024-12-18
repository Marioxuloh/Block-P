package MetricsServer

import (
	client "Block-P/pkg/client"
	metricsClient "Block-P/pkg/client/metrics"
	model "Block-P/pkg/models"
	modelMetrics "Block-P/pkg/models/metrics"
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

func (s *metricsServer) RequestMetricsFromNode(ctx context.Context, req *pb.MetricsRequestTrigger) (*pb.Ack, error) {
	log.Printf("Server: on RequestMetricsFromNode service, received MetricsRequestTrigger from nodeAddress: %v name: %v nodeID: %v", req.NodeAddress, req.Name, req.Id)
	response := &pb.Ack{
		Ack: string("success"),
	}
	go metricsClient.RunNodeMetrics(req.NodeAddress, req.Name, req.Id, 5, 13*time.Second, 15)
	return response, nil
}

func (s *metricsServer) RequestMetrics(req *pb.MetricsRequest, stream pb.MetricService_RequestMetricsServer) error {

	log.Printf("Server: on RequestMetrics service, received  MetricsRequest from nodeID: %v", req.Id)

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

			//GetAddons(), devuelve un map con nombre y respuesta directamente
			//desde metric.go, se llamara al modelo para que acceda a los archivos.pb y nos devuelva las rutas de los scripts paraejecutalosn metric.go
			//el odelodevolveraun con nombre y ruta al script

			metricsAddons, err := modelMetrics.GetAddons()
			if err != nil {
				log.Printf("Error getting Addons: %v", err)
				metricsAddons = nil
			}

			response := &pb.Data{
				Id:      model.GlobalConfig.ID,
				Metrics: metricsAddons,
			}

			if err := stream.Send(response); err != nil {
				log.Printf("Server: on RequestMetrics service, Error sending response for metric %s: %v", metricsAddons, err)
				err := client.Client() //este proceso que es el encargado de hablar con el server invocara un nuevo cliente para establecer otra vez conexion desde 0
				if err != nil {
					log.Printf("Server: Client error %v", err) //si no consigue establecer conexion
					return err
				}
				return err //si consigue establecer conexion tambien acaba
			} else {
				log.Printf("Server: on RequestMetrics service, send response data from %v to master: %v", model.GlobalConfig.Name, response)
				time.Sleep(time.Second) //se envian metricas cada 1 de segundo
			}
		}
	}
}

// ConcatenarMapas concatena dos mapas y devuelve un nuevo mapa
func concatMaps(map1, map2 map[string]string) map[string]string {
	resultado := make(map[string]string)

	// Agregar elementos del primer mapa
	for key, value := range map1 {
		resultado[key] = value
	}

	// Agregar elementos del segundo mapa
	for key, value := range map2 {
		resultado[key] = value
	}

	return resultado
}
