package helpers

import (
	"github.com/lab5e/go-spanuserapi/v4"
	"github.com/lab5e/spancli/pkg/global"
)

// NewSpanAPIClient creates a new SpanAPI client based on the command line options and/or
// defaults.
func NewUserAPIClient() *spanuserapi.APIClient {

	config := spanuserapi.NewConfiguration()
	config.Debug = global.Options.Debug
	if global.Options.OverrideEndpoint != "" {
		config.Servers = spanuserapi.ServerConfigurations{
			spanuserapi.ServerConfiguration{URL: global.Options.OverrideEndpoint},
		}
	}
	credentials := ReadCredentials()
	if credentials != "" {
		config.AddDefaultHeader("Authorization", "TOKEN "+credentials)
	}

	client := spanuserapi.NewAPIClient(config)
	return client
}
