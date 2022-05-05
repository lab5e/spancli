package firmware

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type resetError struct {
	ID commonopt.CollectionAndDevice
}

func (c *resetError) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	_, res, err := client.FotaApi.ClearFirmwareError(ctx, c.ID.CollectionID, c.ID.DeviceID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("Firmware errors on device %s in collection %s cleared\n", c.ID.CollectionID, c.ID.DeviceID)
	return nil
}
