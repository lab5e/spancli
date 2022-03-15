package firmware

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteFirmware struct {
	ID      commonopt.Collection
	ImageID string `long:"image-id" description:"firmware image id"`
	Prompt  commonopt.NoPrompt
}

func (c *deleteFirmware) Execute([]string) error {
	if !c.Prompt.Check() {
		return fmt.Errorf("user aborted delete")
	}
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	image, res, err := client.FotaApi.DeleteFirmware(ctx, c.ID.CollectionID, c.ImageID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("Deleted firmware image %s\n", image.GetImageId())
	return nil
}
