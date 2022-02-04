package gke

import (
	"context"

	"github.com/VixsTy/gke-node-scaler/models"
	"github.com/VixsTy/gke-node-scaler/scaler"
	"google.golang.org/api/container/v1"
)

type GkeScalerService struct {
	name             string
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
	name := "projects/" + event.ProjectId + "/locations/" + event.Zone + "/clusters/" + event.ClusterID + "/nodePools/" + event.NodePoolID
	s.name = name
}

func (s *GkeScalerService) setAutoscaling(enable bool, maxNodeCount int64) error {
	var nodePoolAutoscalingRequest *container.SetNodePoolAutoscalingRequest
	if !enable {
		nodePoolAutoscalingRequest = &container.SetNodePoolAutoscalingRequest{
			Autoscaling: &container.NodePoolAutoscaling{
				Enabled: false,
			},
		}
	} else {
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

	_, err := call.Do()

	return err
}

func (s *GkeScalerService) setNodeSize(nodeCount int64) error {
	nodePoolSizeRequest := container.SetNodePoolSizeRequest{
		Name:      s.name,
		NodeCount: nodeCount,
	}
	call := s.nodePoolsService.SetSize(
		s.name,
		&nodePoolSizeRequest,
	)

	_, err := call.Do()
	if err != nil {
		return err
	}

	return err
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
