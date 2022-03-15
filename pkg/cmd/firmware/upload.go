package firmware

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type uploadFirmware struct {
	ID        commonopt.Collection
	ImageFile string `long:"image" description:"firmware image binary"`
	Version   string `long:"version" description:"firmware image version" default:"1.0.0"`
}

func (c *uploadFirmware) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	bytes, err := os.ReadFile(c.ImageFile)
	if err != nil {
		return err
	}
	if len(bytes) > 1024*1024*4 {
		return fmt.Errorf("image cannot exceed 4MB in size. %s is %2.1f MB", c.ImageFile, (float64(len(bytes)) / (1024.0 * 1024.0)))
	}
	fmt.Printf("Uploading %d bytes from %s...\n", len(bytes), c.ImageFile)

	imageString := base64.StdEncoding.EncodeToString(bytes)
	fwimg, res, err := client.FotaApi.CreateFirmware(ctx, c.ID.CollectionID).Body(
		spanapi.CreateFirmwareRequest{
			Image:    spanapi.PtrString(imageString),
			Version:  spanapi.PtrString(c.Version),
			Filename: spanapi.PtrString(path.Base(c.ImageFile)),
		},
	).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}
	fmt.Printf("Created firmware image with ID %s, SHA256 %s, version %s (%d bytes)\n", fwimg.GetImageId(), fwimg.GetSha256(), fwimg.GetVersion(), fwimg.GetLength())
	return nil
}
