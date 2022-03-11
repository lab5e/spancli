package commonopt

// CollectionAndDevice is a command line parameter prebaked option set.
type CollectionAndDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" env:"SPAN_DEVICE_ID" description:"device id" required:"yes"`
}
