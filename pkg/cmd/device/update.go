package device

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

type updateDevice struct {
	CollectionID     string   `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	NewCollectionID  string   `long:"new-collection-id" description:"Span collection ID you want to move device to"`
	DeviceID         string   `long:"device-id" description:"device id" required:"yes"`
	Name             string   `long:"name" description:"device name"`
	IMSI             string   `long:"imsi" description:"IMSI of device SIM"`
	IMEI             string   `long:"imei" description:"IMEI of device"`
	Tags             []string `long:"tag" description:"set tag value [name:value]"`
	FirmwareTargetID string   `long:"firmware-target-id" description:"set the target firmware id"`
}

func (r *updateDevice) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	device, res, err := client.DevicesApi.RetrieveDevice(ctx, r.CollectionID, r.DeviceID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}
	if r.IMSI != "" {
		device.SetImsi(r.IMSI)
	}
	if r.IMEI != "" {
		device.SetImei(r.IMEI)
	}
	if r.Name != "" {
		r.Tags = append(r.Tags, fmt.Sprintf(`name:"%s"`, r.Name))
	}
	if r.FirmwareTargetID != "" {
		device.Firmware.SetTargetFirmwareId(r.FirmwareTargetID)
	}

	var newCollectionID *string = nil
	if r.NewCollectionID != "" {
		newCollectionID = &r.NewCollectionID
	}

	var firmwareMetadata *spanapi.FirmwareMetadata = nil
	if r.FirmwareTargetID != "" {
		firmwareMetadata = &spanapi.FirmwareMetadata{
			TargetFirmwareId: &r.FirmwareTargetID,
		}
	}

	deviceUpdated, res, err := client.DevicesApi.UpdateDevice(ctx, r.CollectionID, r.DeviceID).Body(spanapi.UpdateDeviceRequest{
		CollectionId: newCollectionID,
		Imsi:         device.Imsi,
		Imei:         device.Imei,
		Tags:         helpers.TagMerge(device.Tags, r.Tags),
		Firmware:     firmwareMetadata,
	}).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("updated device %s\n", *deviceUpdated.DeviceId)
	return nil
}
