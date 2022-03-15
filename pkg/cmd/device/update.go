package device

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type updateDevice struct {
	ID               commonopt.CollectionAndDevice
	NewCollectionID  string `long:"new-collection-id" description:"Span collection ID you want to move device to"`
	Name             string `long:"name" description:"device name"`
	IMSI             string `long:"imsi" description:"IMSI of device SIM"`
	IMEI             string `long:"imei" description:"IMEI of device"`
	Tags             commonopt.Tags
	FirmwareTargetID string `long:"firmware-target-id" description:"set the target firmware id"`
}

func (r *updateDevice) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	update := spanapi.UpdateDeviceRequest{}
	if r.IMSI != "" {
		update.Config = &spanapi.DeviceConfig{
			Ciot: &spanapi.CellularIoTConfig{
				Imsi: &r.IMSI,
			},
		}
	}
	if r.IMEI != "" {
		if update.Config == nil {
			update.Config = &spanapi.DeviceConfig{}
		}
		if update.Config.Ciot == nil {
			update.Config.Ciot = &spanapi.CellularIoTConfig{}
		}
		update.Config.Ciot.Imei = &r.IMEI
	}
	update.Tags = r.Tags.AsMap()
	if r.Name != "" {
		m := *update.Tags
		m["name"] = r.Name
		update.Tags = &m
	}

	if r.NewCollectionID != "" {
		update.CollectionId = spanapi.PtrString(r.NewCollectionID)
	}

	if r.FirmwareTargetID != "" {
		update.Firmware = &spanapi.FirmwareMetadata{
			TargetFirmwareId: spanapi.PtrString(r.FirmwareTargetID),
		}
	}

	deviceUpdated, res, err := client.DevicesApi.UpdateDevice(ctx, r.ID.CollectionID, r.ID.DeviceID).Body(update).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("updated device %s\n", *deviceUpdated.DeviceId)
	return nil
}
