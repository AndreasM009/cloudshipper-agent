package fake

import (
	"fmt"

	base "github.com/andreasM009/cloudshipper-agent/pkg/controller/events"
)

// EventForwarder event forwarder
type EventForwarder struct {
}

// NewFakeEventForwarder new instance
func NewFakeEventForwarder() base.EventForwarder {
	return &EventForwarder{}
}

// ForwardCommandEvent prints event
func (forwarder *EventForwarder) ForwardCommandEvent(evt *base.CommandEvent) {

	var formatstr = `
	Deploymentname:		%s
	DeploymentID:		%s
	JobName:			%s
	JobDisplayname:		%s
	CommandName:		%s
	CommandDisplayname:	%s
	Logs:
		%s
	`

	logs := ""

	for _, l := range evt.Logs {
		if l.IsError {
			logs += fmt.Sprintf("[Error] %s\n", l.Message)
		} else if l.IsWarning {
			logs += fmt.Sprintf("[Warning] %s\n", l.Message)
		} else {
			logs += fmt.Sprintf("[Info] %s\n", l.Message)
		}
	}

	str := fmt.Sprintf(formatstr, evt.DeploymentName, evt.DeploymentID, evt.JobName, evt.JobDisplayName, evt.CommandName, evt.CommandDisplayName, logs)
	fmt.Println(str)
	fmt.Println("")
}

// ForwardDeploymentEvent print
func (forwarder *EventForwarder) ForwardDeploymentEvent(evt *base.DeploymentEvent) {

}
