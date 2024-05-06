package commonopt

// CollectionAndGateway is a command line parameter prebaked option set for a (required) collection ID and a
// (required) gateway ID
type CollectionAndGateway struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	GatewayID    string `long:"gateway-id" env:"SPAN_GATEWAY_ID" description:"gateway id" required:"yes"`
}

// CollectionAndOptionalDevice is a command line parameter prebaked option set for a (required) collection ID
// with an optional device ID
type CollectionAndDeviceOrGateway struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" env:"SPAN_DEVICE_ID" description:"device id" required:"no"`
	GatewayID    string `long:"gateway-id" env:"SPAN_GATEWAY_ID" description:"gateway id" required:"no"`
}

func (c *CollectionAndDeviceOrGateway) Valid() bool {
	if c.DeviceID == "" && c.GatewayID == "" {
		return false
	}
	if c.DeviceID != "" && c.GatewayID != "" {
		return false
	}
	return true
}
