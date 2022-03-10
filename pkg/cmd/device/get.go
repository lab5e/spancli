package device

import (
	"encoding/json"
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
)

type getDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id" required:"yes"`
}

func (r *getDevice) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	device, res, err := client.DevicesApi.RetrieveDevice(ctx, r.CollectionID, r.DeviceID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	jsonData, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}
