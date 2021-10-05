package main

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
