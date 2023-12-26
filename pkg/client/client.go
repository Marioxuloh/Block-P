package Client

import (
	metrics "Block-P/pkg/client/metrics"
	models "Block-P/pkg/models"
	"log"
)

func Client() error {

	//if models.GlobalConfig.Name != "master" {

	for _, Node := range models.GlobalConfig.Nodes {
		if Node.Name == "master" {
			err := metrics.MetricsRequestFromNodeToMaster(Node.Addr, models.GlobalConfig.FullAddress, models.GlobalConfig.Name, models.GlobalConfig.ID)
			if err != nil {
				log.Printf("Client: could not MetricsRequestFromNodeToMaster error %v", err)
				return err
			}
		}
	}

	//}
	return nil
}
