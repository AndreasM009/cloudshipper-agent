package events

// EventForwarder forwards events
type EventForwarder interface {
	ForwardDeploymentEvent(evt *DeploymentEvent)
	ForwardCommandEvent(evt *CommandEvent)
}
