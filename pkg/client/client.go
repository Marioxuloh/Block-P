package Client

import (
	metrics "Block-P/pkg/client/metrics"
	models "Block-P/pkg/models"
	"log"
	"time"
)

func Client() error {
	for {
		for _, Node := range models.GlobalConfig.Nodes {
			if Node.Name == "master" {
				err := metrics.MetricsRequestFromNodeToMaster(Node.Addr, models.GlobalConfig.FullAddress, models.GlobalConfig.Name, models.GlobalConfig.ID)
				if err != nil {
					log.Printf("Client: could not MetricsRequestFromNodeToMaster error %v", err)
					time.Sleep(13 * time.Second) //cada 13 segundos si el mensaje no se ha entregado con exito, de forma infinita
				} else {
					return nil //si el mensaje llega exitosamente se sale del bucle
				}
			}
		}
	}
}
