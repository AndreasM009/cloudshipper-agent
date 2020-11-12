package events

import (
	"time"
)

// Event base
type Event struct {
	DefinitionID   string    `json:"definitionId"`
	DeploymentName string    `json:"deploymentName"`
	DeploymentID   string    `json:"deplyomentId"`
	EventName      string    `json:"eventName"`
	TenantID       string    `json:"tenantId"`
	Timestamp      time.Time `json:"timestamp"`
}
