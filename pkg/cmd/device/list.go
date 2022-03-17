package device

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listDevices struct {
	ID     commonopt.Collection
	Format commonopt.ListFormat
}

func (r *listDevices) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	resp, res, err := client.DevicesApi.ListDevices(ctx, r.ID.CollectionID).Execute()
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
	t.SetTitle("Devices in %s", r.ID.CollectionID)
	t.AppendHeader(table.Row{
		"ID",
		"Config",
		"Network",
		"IP",
		"Last update",
		"Tags"})

	for _, device := range resp.Devices {
		t.AppendRow(table.Row{
			device.GetDeviceId(),
			getConfig(device),
			getMetadata(device),
			getIP(device),
			getLastUpdate(device),
			helpers.TagsToString(device.GetTags()),
		})
	}
	helpers.RenderTable(t, r.Format.Format)
	return nil
}
func getIP(d spanapi.Device) string {
	if d.Metadata.Ciot != nil {
		return d.Metadata.Ciot.GetAllocatedIp()
	}
	if d.Metadata.Inet != nil {
		return d.Metadata.Inet.GetRemoteAddress()
	}
	return "-"
}
func getLastUpdate(d spanapi.Device) string {
	if d.Metadata.Ciot != nil {
		return d.Metadata.Ciot.GetAllocatedAt()
	}
	if d.Metadata.Inet != nil {
		return d.Metadata.Inet.GetLastUpdate()
	}
	return "-"
}

func getConfig(d spanapi.Device) string {
	if d.Config == nil {
		return ""
	}
	if d.Config.HasCiot() && d.Config.Ciot.GetImsi() != "" {
		return fmt.Sprintf("imsi:%s imei:%s", d.Config.Ciot.GetImsi(), d.Config.Ciot.GetImei())
	}
	// No config for inet devices
	return ""
}

func getMetadata(d spanapi.Device) string {
	if !d.HasMetadata() {
		return "(no metadata)"
	}
	if d.Metadata.Ciot != nil {
		operator := d.Metadata.SimOperator
		if operator == nil {
			operator = &spanapi.NetworkOperator{}
		}
		return fmt.Sprintf("cell:%s  operator:%s  country:%s", d.Metadata.Ciot.GetCellId(), operator.GetNetwork(), operator.GetCountry())
	}
	if d.Metadata.Inet != nil {
		return fmt.Sprintf("certificate s/n:%s", d.Metadata.Inet.GetCertificateSerial())
	}
	return ""
}
