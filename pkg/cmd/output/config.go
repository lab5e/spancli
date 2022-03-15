package output

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

// configHelp shows all the configuration options. It's basically just a help text to get an overview over
// the available configuration options for the if
type configHelp struct {
}

func (c *configHelp) Execute([]string) error {
	f := commonopt.ListFormat{
		Format:          "text",
		NoColor:         false,
		PageSize:        0,
		MaxPayloadWdith: 80,
		NumericDate:     false,
	}
	fmt.Println()
	fmt.Println("Webhook (webhook) - Forward payloads to a webhook")
	t := helpers.NewTableOutput(f)
	t.AppendHeader(table.Row{"name", "description"})
	t.AppendRow(table.Row{"url", "URL for webhook"})
	t.AppendRow(table.Row{"basicAuthUser", "Username"})
	t.AppendRow(table.Row{"basicAuthPass", "Password"})
	t.AppendRow(table.Row{"customHeaderName", "Custom header name for request"})
	t.AppendRow(table.Row{"customHeaderValue", "Custom header value for request"})
	helpers.RenderTable(t, "text")

	fmt.Println("UDP (udpout) - Forward payload to an UDP server")
	t = helpers.NewTableOutput(f)
	t.AppendHeader(table.Row{"name", "description"})
	t.AppendRow(table.Row{"host", "(ip/DNS name) Host name for server"})
	t.AppendRow(table.Row{"port", "(numeric) Port number for server"})
	helpers.RenderTable(t, "text")

	fmt.Println("IFTTT (ifttt) - Forward payloads to the IFTTT service")
	t = helpers.NewTableOutput(f)
	t.AppendHeader(table.Row{"name", "description"})
	t.AppendRow(table.Row{"key", "IFTTT secret key"})
	t.AppendRow(table.Row{"event", "Event name"})
	t.AppendRow(table.Row{"asIsPayload", "Deliver payload as is or as JSON"})
	helpers.RenderTable(t, "text")

	fmt.Println("MQTT Client (mqttclient) - Forward payloads to a MQTT broker somewhere on the Internet")
	t = helpers.NewTableOutput(f)
	t.AppendHeader(table.Row{"name", "description"})
	t.AppendRow(table.Row{"endpoint", "MQTT broker endpoint"})
	t.AppendRow(table.Row{"disableCertCheck", "Disable certificate checks"})
	t.AppendRow(table.Row{"username", "MQTT broker username"})
	t.AppendRow(table.Row{"password", "MQTT broker password"})
	t.AppendRow(table.Row{"clientId", "MQTT broker clientid"})
	t.AppendRow(table.Row{"topicName", "Name of topic to publish to"})
	helpers.RenderTable(t, "text")

	fmt.Println("MQTT Broker (mqttbroker) - Publish via the built-in MQTT broker")
	t = helpers.NewTableOutput(f)
	t.AppendHeader(table.Row{"name", "description"})
	t.AppendRow(table.Row{"topicTemple", "Template for topic"})
	t.AppendRow(table.Row{"payloadFormat", "json|binary - format for payload"})
	t.AppendRow(table.Row{"payloadTemplate", "Template for payload"})
	helpers.RenderTable(t, "text")

	return nil
}
