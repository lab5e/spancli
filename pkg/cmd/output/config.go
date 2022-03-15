package output

import (
	"fmt"
)

type configHelp struct {
}

func (c *configHelp) Execute([]string) error {
	fmt.Println("Webhook")
	fmt.Println()
	fmt.Println("               url: URL for webhook")
	fmt.Println("     basicAuthUser: Username")
	fmt.Println("     basicAuthPass: Password")
	fmt.Println("  customHeaderName: Custom header name for request")
	fmt.Println(" customHeaderValue: Custom header value for request")
	fmt.Println()
	fmt.Println("UDP")
	fmt.Println("             host: Host name for server")
	fmt.Println("             port: Port for server")
	fmt.Println()
	fmt.Println("IFTTT")
	fmt.Println("              key: IFTTT secret key")
	fmt.Println("            event: Event name")
	fmt.Println("      asIsPayload: Deliver payload as is or as JSON")
	fmt.Println()
	fmt.Println("MQTT Client")
	fmt.Println("         endpoint: MQTT broker endpoint")
	fmt.Println(" disableCertCheck: Disable certificate checks")
	fmt.Println("         username: MQTT broker username")
	fmt.Println("         password: MQTT broker password")
	fmt.Println("         clientId: MQTT broker clientid")
	fmt.Println("        topicName: Name of topic to publish to")
	fmt.Println()
	fmt.Println("MQTT Broker")
	fmt.Println("     topicTemple: Template for topic")
	fmt.Println("   payloadFormat: json|binary - format for payload")
	fmt.Println(" payloadTemplate: Template for payload")
	fmt.Println()
	return nil
}
