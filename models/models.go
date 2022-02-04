package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type ScalerMessage struct {
	ProjectId    string `json:"project_id"`
	Zone         string `json:"zone"`
	ClusterID    string `json:"clusterId"`
	NodePoolID   string `json:"nodePoolId"`
	NodeCount    int64  `json:"nodeCount"`
	MaxNodeCount int64  `json:"maxNodeCount"`
}

func (m ScalerMessage) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProjectId, validation.Required, is.LowerCase, validation.Length(1, 100)),
		validation.Field(&m.Zone, validation.Required, is.LowerCase, validation.Length(1, 100)),
		validation.Field(&m.ClusterID, validation.Required, is.LowerCase, validation.Length(1, 100)),
		validation.Field(&m.NodePoolID, validation.Required, is.LowerCase, validation.Length(1, 100)),
		validation.Field(&m.NodeCount, validation.Required, validation.Min(0)),
		validation.Field(&m.MaxNodeCount, validation.Required.When(m.NodeCount > 0), validation.Min(1)),
	)
}
