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
		return helpers.APIError(res, err)
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
	ret := ""
	if d.Config.HasCiot() {
		ret += fmt.Sprintf("[ciot imsi:%s imei:%s] ", d.Config.Ciot.GetImsi(), d.Config.Ciot.GetImei())
	}
	if d.Config.HasInet() {
		ret += "[inet] "
	}
	if d.Config.HasGateway() {
		for id := range *d.Config.Gateway {
			ret += fmt.Sprintf("[gw %s] ", id)
		}
	}
	return ret
}

func getMetadata(d spanapi.Device) string {
	if !d.HasMetadata() {
		return "(no metadata)"
	}
	if d.Metadata.Ciot != nil {
		return fmt.Sprintf("cell:%s  network:%s  country:%s", d.Metadata.Ciot.GetCellId(), d.Metadata.Ciot.GetNetwork(), d.Metadata.Ciot.GetCountry())
	}
	if d.Metadata.Inet != nil {
		return fmt.Sprintf("certificate s/n:%s", d.Metadata.Inet.GetCertificateSerial())
	}
	return ""
}
