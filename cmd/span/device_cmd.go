package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/go-spanapi/v4"
)

type deviceCmd struct {
	Add    addDevice    `command:"add" description:"create device"`
	Get    getDevice    `command:"get" description:"get device"`
	Update updateDevice `command:"update" description:"update device"`
	List   listDevices  `command:"list" alias:"ls" description:"list devices"`
	Send   sendDevice   `command:"send" description:"send downstream message"`
	Delete deleteDevice `command:"delete" alias:"del" description:"delete device"`
}

type addDevice struct {
	CollectionID     string   `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	Name             string   `long:"name" description:"device name"`
	IMSI             string   `long:"imsi" description:"IMSI of device SIM" required:"yes"`
	IMEI             string   `long:"imei" description:"IMEI of device" required:"yes"`
	Tags             []string `long:"tag" description:"set tag value [name:value]"`
	FirmwareTargetID string   `long:"firmware-target-id" description:"set the target firmware id"`
}

type getDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id" required:"yes"`
}

type updateDevice struct {
	CollectionID     string   `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	NewCollectionID  string   `long:"new-collection-id" description:"Span collection ID you want to move device to"`
	DeviceID         string   `long:"device-id" description:"device id" required:"yes"`
	Name             string   `long:"name" description:"device name"`
	IMSI             string   `long:"imsi" description:"IMSI of device SIM"`
	IMEI             string   `long:"imei" description:"IMEI of device"`
	Tags             []string `long:"tag" description:"set tag value [name:value]"`
	FirmwareTargetID string   `long:"firmware-target-id" description:"set the target firmware id"`
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
	Transport    string `long:"transport" choice:"udp" choice:"udp-pull" choice:"coap" choice:"coap-pull" description:"transport" required:"yes"` //nolint (choice tags confusess linter)
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
		Tags:         tagMerge(&map[string]string{"name": r.Name}, r.Tags),
		Firmware: &spanapi.FirmwareMetadata{
			TargetFirmwareId: &r.FirmwareTargetID,
		},
	}

	dev, res, err := client.DevicesApi.CreateDevice(ctx, r.CollectionID).Body(device).Execute()
	if err != nil {
		return apiError(res, err)
	}

	fmt.Printf("created device %s\n", *dev.DeviceId)
	return nil
}

func (r *getDevice) Execute([]string) error {
	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	device, res, err := client.DevicesApi.RetrieveDevice(ctx, r.CollectionID, r.DeviceID).Execute()
	if err != nil {
		return apiError(res, err)
	}

	jsonData, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}

func (r *updateDevice) Execute([]string) error {
	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	device, res, err := client.DevicesApi.RetrieveDevice(ctx, r.CollectionID, r.DeviceID).Execute()
	if err != nil {
		return apiError(res, err)
	}
	if r.IMSI != "" {
		device.SetImsi(r.IMSI)
	}
	if r.IMEI != "" {
		device.SetImei(r.IMEI)
	}
	if r.Name != "" {
		r.Tags = append(r.Tags, fmt.Sprintf(`name:"%s"`, r.Name))
	}
	if r.FirmwareTargetID != "" {
		device.Firmware.SetTargetFirmwareId(r.FirmwareTargetID)
	}

	var newCollectionID *string = nil
	if r.NewCollectionID != "" {
		newCollectionID = &r.NewCollectionID
	}

	var firmwareMetadata *spanapi.FirmwareMetadata = nil
	if r.FirmwareTargetID != "" {
		firmwareMetadata = &spanapi.FirmwareMetadata{
			TargetFirmwareId: &r.FirmwareTargetID,
		}
	}

	deviceUpdated, res, err := client.DevicesApi.UpdateDevice(ctx, r.CollectionID, r.DeviceID).Body(spanapi.UpdateDeviceRequest{
		ExistingCollectionId: device.CollectionId,
		DeviceId:             device.DeviceId,
		CollectionId:         newCollectionID,
		Imsi:                 device.Imsi,
		Imei:                 device.Imei,
		Tags:                 tagMerge(device.Tags, r.Tags),
		Firmware:             firmwareMetadata,
	}).Execute()
	if err != nil {
		return apiError(res, err)
	}

	fmt.Printf("updated device %s\n", *deviceUpdated.DeviceId)
	return nil
}

func (r *listDevices) Execute([]string) error {
	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	resp, res, err := client.DevicesApi.ListDevices(ctx, r.CollectionID).Execute()
	if err != nil {
		return apiError(res, err)
	}

	if resp.Devices == nil {
		fmt.Printf("no devices\n")
		return nil
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

func (r *sendDevice) Execute([]string) error {
	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	payload := r.Text
	if !r.IsBase64 {
		payload = base64.StdEncoding.EncodeToString([]byte(r.Text))
	}

	msgResponse, res, err := client.DevicesApi.SendMessage(ctx, r.CollectionID, r.DeviceID).Body(spanapi.SendMessageRequest{
		CollectionId: &r.CollectionID,
		DeviceId:     &r.DeviceID,
		Port:         &r.Port,
		Payload:      &payload,
		Transport:    &r.Transport,
		CoapPath:     &r.CoapPath,
	}).Execute()
	if err != nil {
		return apiError(res, err)
	}

	fmt.Printf("sent %d bytes to device %s\n", *msgResponse.BytesSent, *msgResponse.DeviceId)
	return nil
}

func (r *deleteDevice) Execute([]string) error {
	if !r.YesIAmSure {
		if !verifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	device, res, err := client.DevicesApi.DeleteDevice(ctx, r.CollectionID, r.DeviceID).Execute()
	if err != nil {
		return apiError(res, err)
	}

	fmt.Printf("deleted device %s in collection %s\n", *device.DeviceId, *device.CollectionId)
	return nil
}
