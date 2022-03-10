package global

import "time"

// Global parameters for the commands. This package exists solely to avoid cyclic dependencies. When the
// parser is done the Options variable contains the global parameters for the commands.

type Parameters struct {
	Token            string        `long:"token" env:"SPAN_API_TOKEN" description:"span API token" required:"yes"`
	OverrideEndpoint string        `long:"endpoint" env:"SPAN_API_ENDPOINT" description:"span endpoint override" required:"no"`
	Timeout          time.Duration `long:"timeout" default:"15s" description:"timeout for operation"`
	Debug            bool          `long:"debug" description:"turn on debug output"`
}

var Options *Parameters

func init() {
	Options = new(Parameters)
}
