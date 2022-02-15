package main

import (
	"context"

	function "github.com/VixsTy/gke-node-scaler"
	"github.com/VixsTy/gke-node-scaler/pkg/gke-node-scaler/models"
)

func main() {
	ctx := context.Background()
	m := models.ScalerMessage{
		ProjectId:    "project-id",
		Zone:         "zone",
		ClusterID:    "cluster-id",
		NodePoolID:   "node-pool-id",
		NodeCount:    0,
		MaxNodeCount: 1,
	}
	err := function.ScalerEvent(ctx, m)
	if err != nil {
		panic(err)
	}
}
