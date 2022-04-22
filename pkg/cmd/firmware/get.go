package firmware

import (
	"errors"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type getFirmware struct {
	ID      commonopt.Collection
	ImageID string `long:"image-id" description:"firmware image id" env:"SPAN_IMAGE_ID"`
	Version string `long:"version" description:"firmware version"`
	Format  commonopt.ListFormat
}

func (c *getFirmware) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()
	if c.Version != "" {
		list, res, err := client.FotaApi.ListFirmware(ctx, c.ID.CollectionID).Execute()
		if err != nil {
			return helpers.ApiError(res, err)
		}
		for _, f := range list.Images {
			if f.GetVersion() == c.Version {
				return c.printFirmware(&f)
			}
		}
		return errors.New("unknown version")
	}
	if c.ImageID == "" {
		return errors.New("must specify version or image ID")
	}
	fw, res, err := client.FotaApi.RetrieveFirmware(ctx, c.ID.CollectionID, c.ImageID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	return c.printFirmware(fw)
}

func (c *getFirmware) printFirmware(fw *spanapi.Firmware) error {
	t := helpers.NewTableOutput(c.Format)

	t.SetTitle("Firmware %s", c.ImageID)

	t.AppendHeader(table.Row{"Field", "Value"})

	if err := helpers.DumpToTable(t, fw); err != nil {
		return err
	}

	helpers.RenderTable(t, c.Format.Format)
	return nil
}
