package device

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type addDevice struct {
	ID               commonopt.Collection
	IMSI             string `long:"imsi" description:"IMSI of device SIM"`
	IMEI             string `long:"imei" description:"IMEI of device"`
	Tags             commonopt.Tags
	Eval             bool   `long:"eval" description:"Output device ID as environment variable"`
	FirmwareTargetID string `long:"firmware-target-id" description:"set the target firmware id"`
}

func (r *addDevice) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	device := spanapi.CreateDeviceRequest{
		Tags: r.Tags.AsMap(),
		Firmware: &spanapi.FirmwareMetadata{
			TargetFirmwareId: &r.FirmwareTargetID,
		},
	}
	if r.IMEI != "" && r.IMSI != "" {
		device.Config = &spanapi.DeviceConfig{
			Ciot: &spanapi.CellularIoTConfig{
				Imei: spanapi.PtrString(r.IMEI),
				Imsi: spanapi.PtrString(r.IMSI),
			},
		}
	}
	dev, res, err := client.DevicesApi.CreateDevice(ctx, r.ID.CollectionID).Body(device).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}
	if r.Eval {
		fmt.Printf("export SPAN_DEVICE_ID=%s\n", dev.GetDeviceId())
		return nil
	}
	fmt.Printf("created device %s\n", dev.GetDeviceId())
	return nil
}
