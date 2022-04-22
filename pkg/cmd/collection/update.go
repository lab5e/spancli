package collection

import (
	"errors"
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type updateCollection struct {
	ID   commonopt.Collection
	Tags commonopt.Tags
	//lint:ignore SA5008 Multiple choice tags makes the linter unhappy
	Management      string `long:"firmware-management" description:"firmware management setting" choice:"disabled" choice:"device" choice:"collection"`
	FirmwareVersion string `long:"firmware-version" description:"set the target firmware version"`
}

func (u *updateCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	update := spanapi.UpdateCollectionRequest{
		Tags: u.Tags.AsMap(),
	}
	if u.Management != "" {
		if update.Firmware == nil {
			update.Firmware = &spanapi.CollectionFirmware{}
		}
		switch u.Management {
		case "device":
			update.Firmware.Management = spanapi.FIRMWAREMANAGEMENT_DEVICE.Ptr()
		case "collection":
			update.Firmware.Management = spanapi.FIRMWAREMANAGEMENT_COLLECTION.Ptr()
		case "disabled":
			fallthrough
		default:
			update.Firmware.Management = spanapi.FIRMWAREMANAGEMENT_DISABLED.Ptr()
		}
	}
	if u.FirmwareVersion != "" {
		if update.Firmware == nil {
			update.Firmware = &spanapi.CollectionFirmware{}
		}
		list, res, err := client.FotaApi.ListFirmware(ctx, u.ID.CollectionID).Execute()
		if err != nil {
			return helpers.ApiError(res, err)
		}
		id := ""
		for _, fw := range list.Images {
			if fw.GetVersion() == u.FirmwareVersion {
				id = fw.GetImageId()
				break
			}
		}
		if id == "" {
			return errors.New("unknown version")
		}
		update.Firmware.TargetFirmwareId = spanapi.PtrString(id)
	}
	collectionUpdated, res, err := client.CollectionsApi.UpdateCollection(ctx, u.ID.CollectionID).Body(update).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("updated collection %s\n", *collectionUpdated.CollectionId)
	return nil
}
