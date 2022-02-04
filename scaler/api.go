package scaler

import (
	"context"

	"github.com/VixsTy/gke-node-scaler/models"
)

type ScalerService interface {
	ScaleDown(ctx context.Context, event models.ScalerMessage) error
	ScaleUp(ctx context.Context, event models.ScalerMessage) error
}
