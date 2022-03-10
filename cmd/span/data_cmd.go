package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
	"unicode"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

type dataCmd struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id"`
	Limit        int32  `long:"limit" description:"max number of entries to fetch" default:"30"`
	Start        string `long:"start" description:"start of time range in milliseconds since epoch"`
	End          string `long:"end" description:"end of time range in milliseconds since epoch"`
	Decode       bool   `long:"decode" description:"decode payload"`
	ISODate      bool   `long:"iso-date" description:"format date as ISO-8601 date"`
	JSONOutput   bool   `long:"json" description:"output as JSON"`
	JSONPretty   bool   `long:"pretty" description:"output as indented JSON"`
}

func (r *dataCmd) Execute([]string) error {

	if r.DeviceID == "" {
		return r.listCollectionData()
	}
	return r.listDeviceData()
}

func (r *dataCmd) listDeviceData() error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	req := client.DevicesApi.ListDeviceData(ctx, r.CollectionID, r.DeviceID)
	if r.Start != "" {
		req.Start(r.Start)
	}
	if r.End != "" {
		req.End(r.End)
	}
	req.Limit(r.Limit)
	data, _, err := req.Execute()
	if err != nil {
		return err
	}

	return r.listData(data)
}

func (r *dataCmd) listCollectionData() error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	req := client.CollectionsApi.ListCollectionData(ctx, r.CollectionID)
	if r.Start != "" {
		req.Start(r.Start)
	}
	if r.End != "" {
		req.End(r.End)
	}
	req.Limit(r.Limit)
	data, _, err := req.Execute()
	if err != nil {
		return err
	}

	return r.listData(data)
}

func (r *dataCmd) listData(data *spanapi.ListDataResponse) error {
	if r.JSONOutput {
		fmt.Print("[")

		if r.JSONPretty {
			fmt.Print("\n")
		}

		for n, d := range data.Data {

			var jsonData []byte
			var err error

			if r.JSONPretty {
				jsonData, err = json.MarshalIndent(d, "  ", "    ")
			} else {
				jsonData, err = json.Marshal(d)
			}
			if err != nil {
				return fmt.Errorf("error marshalling to JSON: %v", err)
			}

			if r.JSONPretty {
				fmt.Printf("  ")
			}

			fmt.Print(string(jsonData))

			if n < len(data.Data)-1 {
				fmt.Print(",")
			}
			if r.JSONPretty {
				fmt.Print("\n")
			}

		}
		fmt.Printf("]\n")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, strings.Join([]string{"DeviceID", "Name", "Trans", "Received", "Payload"}, "\t")+"\n")

	for _, d := range data.Data {

		if r.ISODate {
			received, err := strconv.ParseInt(*d.Received, 10, 64)
			if err == nil {
				*d.Received = time.Unix(0, received*int64(time.Millisecond)).Format(time.RFC3339)
			}
		}

		if r.Decode {
			data, err := base64.StdEncoding.DecodeString(*d.Payload)
			if err == nil {

				clean := strings.Map(func(r rune) rune {
					if unicode.IsPrint(r) {
						return r
					}
					return -1
				}, string(data))

				*d.Payload = "'" + clean + "'"
			}
		}
		t := *d.Device.Tags
		fmt.Fprintf(w, strings.Join([]string{
			*d.Device.DeviceId,
			t["name"],
			*d.Transport,
			*d.Received,
			*d.Payload,
		}, "\t")+"\n")
	}

	return w.Flush()

}
