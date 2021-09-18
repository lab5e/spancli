package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/lab5e/spancli/pkg/helpers"
	"github.com/lab5e/spanclient-go/v4"
)

type outputCmd struct {
	Add    addOutput    `command:"add" description:"create output"`
	Remove removeOutput `command:"remove" description:"remove output"`
	Update updateOutput `command:"update" description:"update output"`
	List   listOutputs  `command:"list" alias:"ls" description:"list outputs"`
	Status statusOutput `command:"status" alias:"s" description:"show status"`
	Logs   logsOutput   `command:"logs" alias:"l" description:"show output logs"`
}

type addOutput struct {
	CollectionID      string   `long:"collection-id" short:"c" env:"SPAN_COLLECTION_ID" required:"yes" description:"Span collection ID"`
	OutputID          string   `long:"output-id" short:"o" env:"SPAN_OUTPUT_ID" required:"yes" description:"Span output ID"`
	Type              string   `long:"type" choice:"webhook" choice:"udp" choice:"ifttt" choice:"mqtt-client" choice:"mqtt-broker" description:"output type"`
	Tags              []string `long:"tag" short:"t" description:"Set tag value (name:value)"`
	WebhookURL        string   `long:"webhook-url" description:"URL for endpoint"`
	BasicAuthUser     string   `long:"webhook-basic-auth-user" description:"user name"`
	BasicAuthPass     string   `long:"webhook-basic-auth-pass" description:"password"`
	CustomHeaderName  string   `long:"webhook-custom-header-name" description:"custom header name"`
	CustomHeaderValue string   `long:"webhook-custom-header-value" description:"custom header value"`
	Host              string   `long:"udp-host" description:"Host name"`
	Port              int32    `long:"udp-port" description:"port number"`
	Key               string   `long:"ifttt-key" description:"key"`
	Event             string   `long:"ifttt-event-name" description:"event name"`
	AsIsPayload       bool     `long:"ifttt-as-is-payload" description:"as is payload"`
	Endpoint          string   `long:"mqtt-client-endpoint" description:"MQTT client endpoint"`
	DisableCertCheck  bool     `long:"mqtt-client-disable-cert-check" description:"disable certificate check"`
	TopicName         string   `long:"mqtt-client-topic-name" description:"MQTT client topic name"`
	Username          string   `long:"mqtt-client-username" description:"MQTT client username"`
	Password          string   `long:"mqtt-client-password" description:"MQTT client password"`
	ClientID          string   `long:"mqtt-client-clientid" description:"MQTT client client ID"`
	PayloadType       string   `long:"mqtt-broker-payload-type" choice:"json" choice:"binary" description:"MQTT broker payload"`
	TopicTemplate     string   `long:"mqtt-broker-topic-template" description:"MQTT broker topic template"`
	PayloadTemplate   string   `long:"mqtt-broker-payload-template" description:"MQTT broker payload template"`
}

type updateOutput struct {
	CollectionID string   `long:"collection-id" short:"c" env:"SPAN_COLLECTION_ID" required:"yes" description:"Span collection ID"`
	OutputID     string   `long:"output-id" short:"o" env:"SPAN_OUTPUT_ID" required:"yes" description:"Span output ID"`
	Tags         []string `long:"tag" short:"t" description:"Set tag value (name:value)"`
}

type removeOutput struct {
	CollectionID string `long:"collection-id" short:"c" env:"SPAN_COLLECTION_ID" required:"yes" description:"Span collection ID"`
	OutputID     string `long:"output-id" short:"o" env:"SPAN_OUTPUT_ID" required:"yes" description:"Span output ID"`
	Yes          bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

type listOutputs struct {
	CollectionID string `long:"collection-id" short:"c" env:"SPAN_COLLECTION_ID" required:"yes" description:"Span collection ID"`
}

type statusOutput struct {
	CollectionID string `long:"collection-id" short:"c" env:"SPAN_COLLECTION_ID" required:"yes" description:"Span collection ID"`
	OutputID     string `long:"output-id" short:"o" env:"SPAN_OUTPUT_ID" required:"yes" description:"Span output ID"`
}

type logsOutput struct {
	CollectionID string `long:"collection-id" short:"c" env:"SPAN_COLLECTION_ID" required:"yes" description:"Span collection ID"`
	OutputID     string `long:"output-id" short:"o" env:"SPAN_OUTPUT_ID" required:"yes" description:"Span output ID"`
}

func (a *addOutput) Execute([]string) error {
	tags, err := tagsToMap(a.Tags)
	if err != nil {
		return err
	}
	client := spanclient.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	output, resp, err := client.OutputsApi.CreateOutput(ctx, a.CollectionID, spanclient.Output{
		Type:    spanclient.OutputType(a.Type),
		Enabled: true,
		Config: spanclient.OutputConfig{
			Url:               a.WebhookURL,
			BasicAuthUser:     a.BasicAuthUser,
			BasicAuthPass:     a.BasicAuthPass,
			CustomHeaderName:  a.CustomHeaderName,
			CustomHeaderValue: a.CustomHeaderValue,
			Host:              a.Host,
			Port:              a.Port,
			Key:               a.Key,
			EventName:         a.Event,
			AsIsPayload:       a.AsIsPayload,
			Endpoint:          a.Endpoint,
			DisableCertCheck:  a.DisableCertCheck,
			Username:          a.Username,
			Password:          a.Password,
			ClientId:          a.ClientID,
			TopicName:         a.TopicName,
			// Soon(tm) to be added
			//			PayloadType:       a.PayloadType,
			//			PayloadTemplate:   a.PayloadTemplate,
			//			TopicTemplate:     a.TopicTemplate,
		},
		Tags: tags,
	})
	if err != nil {
		fmt.Printf("Error creating: %+v\n", resp)
		return err
	}
	fmt.Printf("Created output with id '%s'\n", output.OutputId)
	return nil
}

func (u *updateOutput) Execute([]string) error {
	return errors.New("not implemented")
}

func boolToYesNo(enabled bool) string {
	if enabled {
		return "Yes"
	}
	return "No"
}
func (l *listOutputs) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	outputs, _, err := client.OutputsApi.ListOutputs(ctx, l.CollectionID)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintf(w, strings.Join([]string{"OutputID", "Type", "Enabled", "Config", "Tags"}, "\t")+"\n")
	for _, op := range outputs.Outputs {
		fmt.Fprintf(w, strings.Join([]string{
			op.OutputId,
			string(op.Type),
			boolToYesNo(op.Enabled),
			outputConfigToString(op.Type, op.Config),
			tagsToString(op.Tags),
		}, "\t")+"\n")
	}
	return w.Flush()
}

func (e *statusOutput) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	status, _, err := client.OutputsApi.Status(ctx, e.CollectionID, e.OutputID)
	if err != nil {
		return err
	}
	fmt.Printf("Enabled:     %s\n", boolToYesNo(status.Enabled))
	fmt.Printf("Received:    %d\n", status.Received)
	fmt.Printf("Forwarded:   %d\n", status.Forwarded)
	fmt.Printf("Error count: %d\n", status.ErrorCount)
	fmt.Printf("Retransmits: %d\n", status.Retransmits)
	return nil
}

func (l *logsOutput) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	logs, _, err := client.OutputsApi.Logs(ctx, l.CollectionID, l.OutputID)
	if err != nil {
		return err
	}

	if len(logs.Logs) == 0 {
		fmt.Println("No log messages")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)

	fmt.Fprintf(w, strings.Join([]string{"Time", "Message", "Repeated"}, "\t")+"\n")
	for _, logEntry := range logs.Logs {
		fmt.Fprintf(w, strings.Join([]string{
			logEntry.Time,
			logEntry.Message,
			fmt.Sprintf("%d", logEntry.Repeated),
		}, "\t")+"\n")
	}
	return w.Flush()
}

func outputConfigToString(t spanclient.OutputType, config spanclient.OutputConfig) string {
	switch t {
	case spanclient.UDP:
		return fmt.Sprintf("host:%s  port:%d", config.Host, config.Port)
	case spanclient.IFTTT:
		return fmt.Sprintf("key:%s  event:%s  asIsPayload:%t", config.Key, config.EventName, config.AsIsPayload)
	case spanclient.WEBHOOK:
		return fmt.Sprintf("url:%s  user:%s  pass:%s  header:%s  value:%s",
			config.Url, config.BasicAuthUser, config.BasicAuthPass, config.CustomHeaderName, config.CustomHeaderValue)
	case spanclient.MQTT:
		return fmt.Sprintf("endpoint:%s  topic:%s  id:%s  certCheck:%t  user:%s  pass:%s",
			config.Endpoint, config.TopicName, config.ClientId,
			config.DisableCertCheck,
			config.Username, config.Password)

		// TBA
		//	case spanclient.MQTTBROKER:
		//	return fmt.Sprintf("payload:%s  topicTemplate:%s  payloadTemplate:%s", config.PayloadFormat, config.TopicTemplate, config.PayloadTemplate)
	}
	return "unknown"
}
