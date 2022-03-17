package device

import (
	"encoding/json"
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
	// Dump as name.value list by converting to JSON and then back to map[string]any
	buf, err := json.Marshal(device)
	if err != nil {
		return err
	}

	nameValue := map[string]any{}
	if err := json.Unmarshal(buf, &nameValue); err != nil {
		return err
	}
	dumpFields(t, "", nameValue)

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

func fieldName(prefix string, field string) string {
	if prefix == "" {
		return field
	}
	return prefix + "." + field
}

func dumpFields(t table.Writer, prefix string, nameValue map[string]any) {
	for k, v := range nameValue {
		subType, ok := v.(map[string]any)
		if ok {
			dumpFields(t, fieldName(prefix, k), subType)
			continue
		}
		t.AppendRow(table.Row{
			fieldName(prefix, k),
			v,
		})
	}
}
