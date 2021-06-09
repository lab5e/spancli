package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/lab5e/spanclient-go/v4"
)

type listenCmd struct {
	CollectionID   string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID       string `long:"device-id" description:"Span device ID"`
	Pretty         bool   `long:"pretty" description:"pretty-print data"`
	PayloadLogDir  string `long:"payload-log-dir" description:"payload log directory" default:"."`
	PayloadLogName string `long:"payload-log-name" description:"payload log file prefix" default:"payload"`
	LogPayload     bool   `long:"log-payload" description:"log payloads to files suffixed with timestamp in nanoseconds"`
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

		if r.LogPayload {
			r.logPayload(msg.Payload)
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

func (r *listenCmd) logPayload(payload string) {
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		log.Printf("error base64-decoding payload: %v", err)
		return
	}

	filename := fmt.Sprintf("%s/%s-%d", r.PayloadLogDir, r.PayloadLogName, time.Now().UnixNano())
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatalf("Error writing payload to %s: %v", filename, err)
	}
	log.Printf("wrote file: %s", filename)
}
