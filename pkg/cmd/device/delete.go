package device

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteDevice struct {
	ID     commonopt.CollectionAndDevice
	Prompt commonopt.NoPrompt
}

func (r *deleteDevice) Execute([]string) error {
	if !r.Prompt.Check() {
		return fmt.Errorf("user aborted delete")
	}

	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	device, res, err := client.DevicesApi.DeleteDevice(ctx, r.ID.CollectionID, r.ID.DeviceID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("deleted device %s in collection %s\n", *device.DeviceId, *device.CollectionId)
	return nil
}
