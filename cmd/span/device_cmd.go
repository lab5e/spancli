package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/lab5e/spanclient-go/v4"
)

type deviceCmd struct {
	Add    addDevice    `command:"add" description:"create device"`
	Get    getDevice    `command:"get" description:"get device"`
	List   listDevice   `command:"list" description:"list devices"`
	Send   sendDevice   `command:"send" description:"send downstream message"`
	Delete deleteDevice `command:"delete" description:"delete device"`
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

type listDevice struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
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

func init() {
	parser.AddCommand("device", "device management commands", "device management commands", &deviceCmd{})
}

func (r *addDevice) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, _ := spanContext()

	device, _, err := client.DevicesApi.CreateDevice(ctx, r.CollectionID, spanclient.Device{
		CollectionId: r.CollectionID,
		Imsi:         r.IMSI,
		Imei:         r.IMEI,
		Tags:         map[string]string{"name": r.Name},
	})
	if err != nil {
		return err
	}

	fmt.Printf("created device with id '%s'\n", device.DeviceId)
	return nil
}

func (r *getDevice) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, _ := spanContext()
	device, _, err := client.DevicesApi.RetrieveDevice(ctx, r.CollectionID, r.DeviceID)
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
	client := spanclient.NewAPIClient(clientConfig())
	ctx, _ := spanContext()

	devices, _, err := client.DevicesApi.ListDevices(ctx, r.CollectionID)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, strings.Join([]string{"DeviceID", "Name", "IMSI", "IMEI", "IP", "At", "Cell", "FW", "State"}, "\t")+"\n")
	for _, dev := range devices.Devices {

		dev.Tags["name"] = truncateString(dev.Tags["name"], 25)

		fmt.Fprintf(w, strings.Join([]string{
			dev.DeviceId,
			dev.Tags["name"],
			dev.Imsi,
			dev.Imei,
			dev.Network.AllocatedIp,
			dev.Network.AllocatedAt,
			dev.Network.CellId,
			dev.Firmware.FirmwareVersion,
			dev.Firmware.State,
		}, "\t")+"\n")
	}
	return w.Flush()
}

func (r *sendDevice) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, _ := spanContext()

	payload := r.Text
	if !r.IsBase64 {
		payload = base64.StdEncoding.EncodeToString([]byte(r.Text))
	}

	_, _, err := client.DevicesApi.SendMessage(ctx, r.CollectionID, r.DeviceID, spanclient.SendMessageRequest{
		CollectionId: r.CollectionID,
		DeviceId:     r.DeviceID,
		Port:         r.Port,
		Payload:      payload,
		Transport:    r.Transport,
		CoapPath:     r.CoapPath,
	})
	return err
}

func (r *deleteDevice) Execute([]string) error {
	if !r.YesIAmSure {
		if !verifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client := spanclient.NewAPIClient(clientConfig())
	ctx, _ := spanContext()

	_, _, err := client.DevicesApi.DeleteDevice(ctx, r.CollectionID, r.DeviceID)
	if err != nil {
		return err
	}

	fmt.Printf("deleted device '%s'\n", r.DeviceID)
	return nil
}
