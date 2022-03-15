package device

import (
	"encoding/json"
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type getDevice struct {
	ID commonopt.CollectionAndDevice
}

func (r *getDevice) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	device, res, err := client.DevicesApi.RetrieveDevice(ctx, r.ID.CollectionID, r.ID.DeviceID).Execute()
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
