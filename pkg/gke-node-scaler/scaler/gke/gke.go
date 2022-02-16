package gke

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/VixsTy/gke-node-scaler/pkg/gke-node-scaler/models"
	"github.com/VixsTy/gke-node-scaler/pkg/gke-node-scaler/scaler"
	container "google.golang.org/api/container/v1beta1"
)

type GkeScalerService struct {
	name             string
	projectID        string
	containerService *container.Service
	nodePoolsService *container.ProjectsLocationsClustersNodePoolsService
}

func NewScalerService() scaler.ScalerService {
	return &GkeScalerService{}
}

func (s *GkeScalerService) newContainerService(ctx context.Context) error {
	containerService, err := container.NewService(ctx)
	if err != nil {
		return err
	}
	s.containerService = containerService
	s.nodePoolsService = container.NewProjectsLocationsClustersNodePoolsService(containerService)

	return nil
}

func (s *GkeScalerService) buildName(event models.ScalerMessage) {
	name := "projects/" + event.ProjectID + "/locations/" + event.Zone + "/clusters/" + event.ClusterID + "/nodePools/" + event.NodePoolID
	s.name = name
	s.projectID = event.ProjectID
}

func (s *GkeScalerService) waitForOperationEnding(op *container.Operation) (err error) {
	opFullName := fmt.Sprintf("projects/%s/locations/%s/operations/%s", s.projectID, op.Zone, op.Name)
	i := 0
	for i < 10 && op.Status != "DONE" {
		log.Println("Waiting for operation to complete...", op.Name, s.name)
		time.Sleep(30 * time.Second)
		op, err = container.NewProjectsLocationsOperationsService(s.containerService).Get(opFullName).Do()
		if err != nil {
			return err
		}
		i++
	}
	log.Println("Operation complete !", op.Name, s.name)
	return nil
}

func (s *GkeScalerService) setAutoscaling(enable bool, maxNodeCount int64) error {
	var nodePoolAutoscalingRequest *container.SetNodePoolAutoscalingRequest
	if !enable {
		log.Println("Disabling autoscaling...", s.name)
		nodePoolAutoscalingRequest = &container.SetNodePoolAutoscalingRequest{
			Autoscaling: &container.NodePoolAutoscaling{
				Enabled: false,
			},
		}
	} else {
		log.Println("Enabling autoscaling...", maxNodeCount, s.name)
		nodePoolAutoscalingRequest = &container.SetNodePoolAutoscalingRequest{
			Autoscaling: &container.NodePoolAutoscaling{
				Enabled:      true,
				MaxNodeCount: maxNodeCount,
			},
		}
	}

	call := s.nodePoolsService.SetAutoscaling(
		s.name,
		nodePoolAutoscalingRequest,
	)

	op, err := call.Do()
	if err != nil {
		return err
	}

	return s.waitForOperationEnding(op)
}

func (s *GkeScalerService) setNodeSize(nodeCount int64) error {
	log.Println("Setting node size...", nodeCount, s.name)
	nodePoolSizeRequest := container.SetNodePoolSizeRequest{
		Name:      s.name,
		NodeCount: nodeCount,
	}
	call := s.nodePoolsService.SetSize(
		s.name,
		&nodePoolSizeRequest,
	)

	op, err := call.Do()
	if err != nil {
		return err
	}

	return s.waitForOperationEnding(op)
}

func (s *GkeScalerService) ScaleDown(ctx context.Context, event models.ScalerMessage) error {
	err := s.newContainerService(ctx)
	if err != nil {
		return err
	}

	s.buildName(event)

	err = s.setAutoscaling(false, 0)
	if err != nil {
		return err
	}

	err = s.setNodeSize(event.NodeCount)
	if err != nil {
		return err
	}

	return nil
}

func (s *GkeScalerService) ScaleUp(ctx context.Context, event models.ScalerMessage) error {
	err := s.newContainerService(ctx)
	if err != nil {
		return err
	}

	s.buildName(event)

	err = s.setAutoscaling(event.MaxNodeCount > 0, event.MaxNodeCount)
	if err != nil {
		return err
	}

	err = s.setNodeSize(event.NodeCount)
	if err != nil {
		return err
	}

	return nil
}
