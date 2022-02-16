package function

import (
	"context"
	"log"

	"github.com/VixsTy/gke-node-scaler/pkg/gke-node-scaler/models"
	"github.com/VixsTy/gke-node-scaler/pkg/gke-node-scaler/scaler/gke"
)

// ScalerEvent consumes a Pub/Sub message.
//nolint:deadcode,unused
// nolint cause it's designed to be used as a function
func ScalerEvent(ctx context.Context, m models.ScalerMessage) error {

	log.Println("ScalerEvent: New event", m)

	log.Println("ScalerEvent: create GKE Scaler service")
	scalerService := gke.NewScalerService()

	var err error
	if m.NodeCount == 0 {
		log.Println("ScalerEvent: Scale down")
		err = scalerService.ScaleDown(ctx, m)
	} else if m.NodeCount > 0 {
		log.Println("ScalerEvent: Scale up")
		err = scalerService.ScaleUp(ctx, m)
	}

	return err
}
