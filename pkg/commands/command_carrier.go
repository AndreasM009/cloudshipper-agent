package commands

// CommandCarrier carries command for streaming it to outproc components
type CommandCarrier struct {
	CarrierForType CommandType `json:"carrierForType"`
	Data           interface{} `json:"data"`
}
