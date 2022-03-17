package device

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type getDevice struct {
	ID     commonopt.CollectionAndDevice
	Format commonopt.ListFormat
}

func (r *getDevice) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	device, res, err := client.DevicesApi.RetrieveDevice(ctx, r.ID.CollectionID, r.ID.DeviceID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	t := helpers.NewTableOutput(r.Format)
	t.SetTitle("Device %s", device.GetDeviceId())

	t.AppendHeader(table.Row{"Field", "Value"})

	if err := helpers.DumpToTable(t, device); err != nil {
		return err
	}

	certs, res, err := client.DevicesApi.DeviceCertificate(ctx, r.ID.CollectionID, r.ID.DeviceID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	// List certificates separately
	for i, v := range certs.Certificates {
		t.AppendRow(table.Row{
			fmt.Sprintf("device.certificate.%d.serial", i),
			v.GetCertificateSerial(),
		})
		t.AppendRow(table.Row{
			fmt.Sprintf("device.certificate.%d.expires", i),
			helpers.DateFormat(v.GetExpires(), r.Format.NumericDate),
		})
	}
	helpers.RenderTable(t, r.Format.Format)
	return nil
}
