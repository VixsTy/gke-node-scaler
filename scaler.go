package function

import (
	"context"

	"github.com/VixsTy/gke-node-scaler/pkg/gke-node-scaler/models"
	"github.com/VixsTy/gke-node-scaler/pkg/gke-node-scaler/scaler/gke"
)

// ScalerEvent consumes a Pub/Sub message.
//nolint:deadcode,unused
// nolint cause it's designed to be used as a function
func ScalerEvent(ctx context.Context, m models.ScalerMessage) error {

	scalerService := gke.NewScalerService()

	var err error
	if m.NodeCount == 0 {
		err = scalerService.ScaleDown(ctx, m)
	} else if m.NodeCount > 0 {
		err = scalerService.ScaleUp(ctx, m)
	}

	return err
}
