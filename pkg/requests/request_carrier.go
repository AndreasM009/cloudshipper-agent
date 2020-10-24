package requests

// RequestCarrier for marshaling
type RequestCarrier struct {
	CarrierForType RequestType `json:"type"`
	Data           interface{} `jason:"data"`
}
