package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/lab5e/spanclient-go/v4"
)

type listenCmd struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"Span device ID"`
	Pretty       bool   `long:"pretty" description:"pretty-print data"`
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

		var jsonBytes []byte
		if r.Pretty {
			jsonBytes, err = json.MarshalIndent(msg, "", "    ")
		} else {
			jsonBytes, err = json.Marshal(msg)
		}
		if err != nil {
			log.Printf("JSON marshalling error: %v", err)
			continue
		}

		fmt.Printf("%s\n", jsonBytes)
	}
}

// openDataStream will open a DataStream from a device if DeviceID is
// specified, otherwise it will return a DataStream from the
// collection.
func (r *listenCmd) openDataStream() (spanclient.DataStream, error) {
	ctx, _ := spanContext()
	if r.DeviceID != "" {
		return spanclient.NewDeviceDataStream(ctx, clientConfig(), r.CollectionID, r.DeviceID)
	}
	return spanclient.NewCollectionDataStream(ctx, clientConfig(), r.CollectionID)
}
