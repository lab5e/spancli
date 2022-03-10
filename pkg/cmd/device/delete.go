package device

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id" required:"yes"`
	YesIAmSure   bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

func (r *deleteDevice) Execute([]string) error {
	if !r.YesIAmSure {
		if !helpers.VerifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	device, res, err := client.DevicesApi.DeleteDevice(ctx, r.CollectionID, r.DeviceID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("deleted device %s in collection %s\n", *device.DeviceId, *device.CollectionId)
	return nil
}
