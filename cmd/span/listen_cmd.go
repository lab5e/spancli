package main

import (
	"fmt"
	"log"

	"github.com/lab5e/spanclient-go"
)

type listenCmd struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" env:"SPAN_DEVICE_ID" description:"Span device ID"`
}

func init() {
	parser.AddCommand("listen", "Listen for messages from Span", "Listen for messages from Span", &listenCmd{})
}

func (r *listenCmd) Execute([]string) error {
	ds, err := r.openDataStream()
	if err != nil {
		return err
	}
	defer ds.Close()

	for {
		msg, err := ds.Recv()
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", msg)
	}
}

// openDataStream will open a DataStream from a device if DeviceID is
// specified, otherwise it will return a DataStream from the
// collection.
func (r *listenCmd) openDataStream() (spanclient.DataStream, error) {
	ctx, _ := spanContext()

	log.Printf("Context: %+v", ctx)

	if r.DeviceID != "" {
		return spanclient.NewDeviceDataStream(ctx, clientConfig(), r.CollectionID, r.DeviceID)
	}
	return spanclient.NewCollectionDataStream(ctx, clientConfig(), r.CollectionID)
}
