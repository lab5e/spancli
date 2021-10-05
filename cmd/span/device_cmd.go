package main

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/go-spanapi/v4"
)

type deviceCmd struct {
	Add    addDevice    `command:"add" description:"create device"`
	Get    getDevice    `command:"get" description:"get device"`
	List   listDevices  `command:"list" alias:"ls" description:"list devices"`
	Send   sendDevice   `command:"send" description:"send downstream message"`
	Delete deleteDevice `command:"delete" alias:"del" description:"delete device"`
}

type addDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	Name         string `long:"name" description:"device name" required:"yes"`
	IMSI         string `long:"imsi" description:"IMSI of device SIM" required:"yes"`
	IMEI         string `long:"imei" description:"IMEI of device" required:"yes"`
}

type getDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id" required:"yes"`
}

type listDevices struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	Format       string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor      bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize     int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

type sendDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id" required:"yes"`
	Port         int32  `long:"port" default:"1234" description:"destination port on device" required:"yes"`
	Transport    string `long:"transport" choice:"udp-push" choice:"udp-pull"  choice:"coap-push" choice:"coap-pull" description:"transport" required:"yes"`
	CoapPath     string `long:"coap-path" description:"CoAP path"`
	Text         string `long:"text" description:"text payload" required:"yes"`
	IsBase64     bool   `long:"base64" description:"indicates that --text is base64 data"`
}

type deleteDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id" required:"yes"`
	YesIAmSure   bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

func (r *addDevice) Execute([]string) error {
	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	device := spanapi.Device{
		CollectionId: &r.CollectionID,
		Imsi:         &r.IMSI,
		Imei:         &r.IMEI,
		Tags:         &map[string]string{},
	}

	if r.Name != "" {
		(*device.Tags)["name"] = r.Name
	}

	dev, res, err := client.DevicesApi.CreateDevice(ctx, r.CollectionID).Body(device).Execute()
	if err != nil {
		return apiError(res, err)
	}

	fmt.Printf("created device %s\n", *dev.DeviceId)
	return nil
}

func (r *listDevices) Execute([]string) error {
	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	resp, res, err := client.DevicesApi.ListDevices(ctx, r.CollectionID).Execute()
	if err != nil {
		return apiError(res, err)
	}

	if r.Format == "json" {
		json, err := json.MarshalIndent(resp.Devices, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := newTableOutput(r.Format, r.NoColor, r.PageSize)
	t.SetTitle("Devices in %s", r.CollectionID)
	t.AppendHeader(table.Row{"DeviceID", "Name", "Last conn", "FW", "IMSI", "IMEI"})

	for _, device := range *resp.Devices {
		// only truncate name if we output as 'text'
		name := device.GetTags()["name"]
		if r.Format == "text" {
			name = truncateString(name, 25)
		}

		allocatedAt := "-"
		if *device.Network.AllocatedAt != "0" {
			allocatedAt = localTimeFormat(*device.Network.AllocatedAt)
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
	renderTable(t, r.Format)

	return nil
}
