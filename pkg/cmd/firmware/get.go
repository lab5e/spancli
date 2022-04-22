package firmware

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type getFirmware struct {
	ID      commonopt.Collection
	ImageID string `long:"image-id" description:"firmware image id" required:"yes" env:"SPAN_IMAGE_ID"`
	Format  commonopt.ListFormat
}

func (c *getFirmware) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	fw, res, err := client.FotaApi.RetrieveFirmware(ctx, c.ID.CollectionID, c.ImageID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	t := helpers.NewTableOutput(c.Format)

	t.SetTitle("Firmware %s", c.ImageID)

	t.AppendHeader(table.Row{"Field", "Value"})

	if err := helpers.DumpToTable(t, fw); err != nil {
		return err
	}

	helpers.RenderTable(t, c.Format.Format)
	return nil
}
