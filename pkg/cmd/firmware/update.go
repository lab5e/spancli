package firmware

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type updateFirmware struct {
	ID              commonopt.Collection
	Tags            commonopt.Tags
	Version         string `long:"version" description:"version number for firmware"`
	ImageID         string `long:"image-id" description:"firmware image id"`
	NewCollectionID string `long:"new-collection-id" description:"new collection id for image"`
}

func (c *updateFirmware) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	updateRequest := spanapi.UpdateFirmwareRequest{}
	updateRequest.Tags = c.Tags.AsMap()

	if c.Version != "" {
		updateRequest.Version = &c.Version
	}
	if c.NewCollectionID != "" {
		updateRequest.CollectionId = &c.NewCollectionID
	}
	fw, res, err := client.FotaApi.UpdateFirmware(ctx, c.ID.CollectionID, c.ImageID).Body(updateRequest).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}
	fmt.Printf("Updated firmware image %s\n", fw.GetImageId())
	return nil
}
