package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type ScalerMessage struct {
	ProjectID    string `json:"project_id"`
	Zone         string `json:"zone"`
	ClusterID    string `json:"cluster_id"`
	NodePoolID   string `json:"node_pool_id"`
	NodeCount    int64  `json:"node_count"`
	MaxNodeCount int64  `json:"max_node_count"`
}

func (m ScalerMessage) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProjectID, validation.Required, is.LowerCase, validation.Length(1, 100)),
		validation.Field(&m.Zone, validation.Required, is.LowerCase, validation.Length(1, 100)),
		validation.Field(&m.ClusterID, validation.Required, is.LowerCase, validation.Length(1, 100)),
		validation.Field(&m.NodePoolID, validation.Required, is.LowerCase, validation.Length(1, 100)),
		validation.Field(&m.NodeCount, validation.Required, validation.Min(0)),
		validation.Field(&m.MaxNodeCount, validation.Required.When(m.NodeCount > 0), validation.When(m.NodeCount > 0, validation.Min(1))),
	)
}
