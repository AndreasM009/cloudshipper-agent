package azure

import "github.com/andreasM009/cloudshipper-agent/pkg/commands"

// AgentAzCommand is an command executed against Azure
type AgentAzCommand struct {
	commands.Command       `json:"command"`
	ServicePrincipalName   string `json:"servicePirncipalName"`
	ServicePrincipalSecret string `json:"servicePrincipalSecret"`
	Tenant                 string `json:"Tenant"`
	Susbcription           string `json:"Subscription"`
}
