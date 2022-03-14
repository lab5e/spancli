package device

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listDevices struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	Format       commonopt.ListFormat
}

func (r *listDevices) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	resp, res, err := client.DevicesApi.ListDevices(ctx, r.CollectionID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	if resp.Devices == nil {
		fmt.Printf("no devices\n")
		return nil
	}

	if r.Format.Format == "json" {
		json, err := json.MarshalIndent(resp.Devices, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := helpers.NewTableOutput(r.Format)
	t.SetTitle("Devices in %s", r.CollectionID)
	t.AppendHeader(table.Row{"DeviceID", "Name", "Last conn", "FW", "IMSI", "IMEI"})

	for _, device := range resp.Devices {
		// only truncate name if we output as 'text'
		name := device.GetTags()["name"]
		if r.Format.Format == "text" {
			name = helpers.EllipsisString(name, 25)
		}

		allocatedAt := "-"
		if *device.Network.AllocatedAt != "0" {
			allocatedAt = helpers.LocalTimeFormat(*device.Network.AllocatedAt)
		}

		fwVersion := "-"
		if *device.Firmware.FirmwareVersion != "" {
			fwVersion = *device.Firmware.FirmwareVersion
		}

		t.AppendRow(table.Row{
			*device.DeviceId,
			name,
			allocatedAt,
			fwVersion,
			*device.Imsi,
			*device.Imei,
		})
	}
	helpers.RenderTable(t, r.Format.Format)

	return nil
}
