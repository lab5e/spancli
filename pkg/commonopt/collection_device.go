package commonopt

// CollectionAndDevice is a command line parameter prebaked option set for a (required) collection ID and a
// (required) device ID
type CollectionAndDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" env:"SPAN_DEVICE_ID" description:"device id" required:"yes"`
}

// CollectionAndOptionalDevice is a command line parameter prebaked option set for a (required) collection ID
// with an optional device ID
type CollectionAndOptionalDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" env:"SPAN_DEVICE_ID" description:"device id" required:"no"`
}
