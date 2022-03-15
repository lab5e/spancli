package firmware

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listFirmware struct {
	ID     commonopt.Collection
	Format commonopt.ListFormat
}

func (c *listFirmware) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	list, res, err := client.FotaApi.ListFirmware(ctx, c.ID.CollectionID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	t := helpers.NewTableOutput(c.Format)
	t.SetTitle("Firmware images for collection %s", c.ID.CollectionID)
	t.AppendHeader(table.Row{
		"Image ID",
		"Created",
		"Version",
		"Filename",
		"Length",
	})
	for _, fw := range list.Images {
		t.AppendRow(table.Row{
			fw.GetImageId(),
			helpers.DateFormat(fw.GetCreated(), c.Format.NumericDate),
			fw.GetVersion(),
			fw.GetFilename(),
			fw.GetLength(),
		})
	}
	t.AppendFooter(table.Row{fmt.Sprintf("%d image(s)", len(list.Images))})
	helpers.RenderTable(t, c.Format.Format)

	return nil
}
