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
	Eval             bool     `long:"eval" description:"Output device ID as environment variable"`
	FirmwareTargetID string   `long:"firmware-target-id" description:"set the target firmware id"`
	GatewayID        string   `long:"gateway-id" description:"configuration for gateway" `
	ConfigParam      []string `long:"config" description:"configuration parameters"`
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
	if len(r.ConfigParam) > 0 && r.GatewayID == "" {
		return fmt.Errorf("must specify gateway ID when configuring device")
	}
	device.Config = &spanapi.DeviceConfig{}
	if r.IMEI != "" && r.IMSI != "" {
		device.Config.Ciot = &spanapi.CellularIoTConfig{
			Imei: spanapi.PtrString(r.IMEI),
			Imsi: spanapi.PtrString(r.IMSI),
		}
	}
	if r.GatewayID != "" {
		device.Config.Gateway = &map[string]spanapi.GatewayDeviceConfig{
			r.GatewayID: {
				GatewayId: spanapi.PtrString(r.GatewayID),
				Params:    helpers.AsMap(r.ConfigParam),
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
