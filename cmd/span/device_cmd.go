package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

type deviceCmd struct {
	Add    addDevice    `command:"add" description:"create device"`
	Get    getDevice    `command:"get" description:"get device"`
	List   listDevice   `command:"list" alias:"ls" description:"list devices"`
	Send   sendDevice   `command:"send" description:"send downstream message"`
	Delete deleteDevice `command:"delete" alias:"del" description:"delete device"`
}

type addDevice struct {
	CollectionID string   `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	IMSI         string   `long:"imsi" description:"IMSI of device SIM" required:"yes"`
	IMEI         string   `long:"imei" description:"IMEI of device" required:"yes"`
	Tags         []string `long:"tag" short:"t" description:"Set tag value (name:value)"`
}

type getDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id" required:"yes"`
}

type listDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
}

type sendDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id" required:"yes"`
	Port         int32  `long:"port" default:"1234" description:"destination port on device" required:"yes"`
	Transport    string `long:"transport" choice:"udp-push" choice:"udp-pull" choice:"coap-push" choice:"coap-pull" description:"transport" required:"yes"` //nolint (choice tags confusess linter)
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
	tags, err := tagsToMap(r.Tags)
	if err != nil {
		return err
	}
	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	device, _, err := client.DevicesApi.CreateDevice(ctx, r.CollectionID).Body(
		spanapi.Device{
			CollectionId: spanapi.PtrString(r.CollectionID),
			Imsi:         spanapi.PtrString(r.IMSI),
			Imei:         spanapi.PtrString(r.IMEI),
			Tags:         &tags,
		}).Execute()
	if err != nil {
		return err
	}

	fmt.Printf("created device with id '%s'\n", *device.DeviceId)
	return nil
}

func (r *getDevice) Execute([]string) error {
	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	device, _, err := client.DevicesApi.RetrieveDevice(ctx, r.CollectionID, r.DeviceID).Execute()
	if err != nil {
		return err
	}
	json, err := json.MarshalIndent(device, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to marshal '%v' to JSON: %v", device, err)
	}
	fmt.Printf("%s\n", json)
	return nil
}

func (r *listDevice) Execute([]string) error {
	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	devices, _, err := client.DevicesApi.ListDevices(ctx, r.CollectionID).Execute()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, strings.Join([]string{"DeviceID", "IMSI", "IMEI", "IP", "At", "Cell", "FW version", "State", "Tags"}, "\t")+"\n")
	for _, dev := range *devices.Devices {
		fmt.Fprintf(w, strings.Join([]string{
			strPtr(dev.DeviceId),
			strPtr(dev.Imsi),
			strPtr(dev.Imei),
			strPtr(dev.Network.AllocatedIp),
			strPtr(dev.Network.AllocatedAt),
			strPtr(dev.Network.CellId),
			strPtr(dev.Firmware.FirmwareVersion),
			strPtr(dev.Firmware.State),
			tagsToString(*dev.Tags),
		}, "\t")+"\n")
	}
	return w.Flush()
}

func (r *sendDevice) Execute([]string) error {
	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	payload := r.Text
	if !r.IsBase64 {
		payload = base64.StdEncoding.EncodeToString([]byte(r.Text))
	}

	_, _, err := client.DevicesApi.SendMessage(ctx, r.CollectionID, r.DeviceID).Body(
		spanapi.SendMessageRequest{
			CollectionId: spanapi.PtrString(r.CollectionID),
			DeviceId:     spanapi.PtrString(r.DeviceID),
			Port:         spanapi.PtrInt32(r.Port),
			Payload:      spanapi.PtrString(payload),
			Transport:    spanapi.PtrString(r.Transport),
			CoapPath:     spanapi.PtrString(r.CoapPath),
		}).Execute()
	return err
}

func (r *deleteDevice) Execute([]string) error {
	if !r.YesIAmSure {
		if !verifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	_, _, err := client.DevicesApi.DeleteDevice(ctx, r.CollectionID, r.DeviceID).Execute()
	if err != nil {
		return err
	}

	fmt.Printf("deleted device '%s'\n", r.DeviceID)
	return nil
}
