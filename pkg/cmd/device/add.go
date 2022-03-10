package device

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

type addDevice struct {
	CollectionID     string   `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	Name             string   `long:"name" description:"device name"`
	IMSI             string   `long:"imsi" description:"IMSI of device SIM" required:"yes"`
	IMEI             string   `long:"imei" description:"IMEI of device" required:"yes"`
	Tags             []string `long:"tag" description:"set tag value [name:value]"`
	FirmwareTargetID string   `long:"firmware-target-id" description:"set the target firmware id"`
}

func (r *addDevice) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	device := spanapi.CreateDeviceRequest{
		Imsi: &r.IMSI,
		Imei: &r.IMEI,
		Tags: helpers.TagMerge(&map[string]string{"name": r.Name}, r.Tags),
		Firmware: &spanapi.FirmwareMetadata{
			TargetFirmwareId: &r.FirmwareTargetID,
		},
	}

	dev, res, err := client.DevicesApi.CreateDevice(ctx, r.CollectionID).Body(device).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("created device %s\n", *dev.DeviceId)
	return nil
}
