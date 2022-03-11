package collection

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

type updateCollection struct {
	CollectionID string   `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	Tags         []string `long:"tag" description:"Set tag value (name:value)"`
	//lint:ignore SA5008 Multiple choice tags makes the linter unhappy
	Management       string `long:"firmware-management" description:"firmware management setting" choice:"disabled" choice:"device" choice:"collection"`
	FirmwareTargetID string `long:"firmware-target-id" description:"set the target firmware id"`
}

func (u *updateCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	update := spanapi.UpdateCollectionRequest{
		Tags: helpers.TagMerge(nil, u.Tags),
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
	if u.FirmwareTargetID != "" {
		if update.Firmware == nil {
			update.Firmware = &spanapi.CollectionFirmware{}
		}
		update.Firmware.TargetFirmwareId = spanapi.PtrString(u.FirmwareTargetID)
	}
	collectionUpdated, res, err := client.CollectionsApi.UpdateCollection(ctx, u.CollectionID).Body(update).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("updated collection %s\n", *collectionUpdated.CollectionId)
	return nil
}
