package global

import "time"

// Global parameters for the commands. This package exists solely to avoid cyclic dependencies. When the
// parser is done the Options variable contains the global parameters for the commands.

// Parameters is the global parameters for all the commands
type Parameters struct {
	Token                string        `long:"token" env:"SPAN_API_TOKEN" description:"span API token"`
	OverrideEndpoint     string        `long:"endpoint" env:"SPAN_API_ENDPOINT" description:"span endpoint override"`
	MQTTOverrideEndpoint string        `long:"mqtt-endpoint" env:"SPAN_MQTT_ENDPOINT" description:"span MQTT endpoint override"`
	Timeout              time.Duration `long:"timeout" default:"15s" description:"timeout for operation"`
	Debug                bool          `long:"debug" description:"turn on debug output"`
}

// Options holds the global options
var Options *Parameters

func init() {
	Options = new(Parameters)
}
