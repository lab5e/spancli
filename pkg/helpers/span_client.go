package helpers

import (
	"context"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/go-spanapi/v4/apitools"
	"github.com/lab5e/spancli/pkg/global"
)

// NewSpanAPIClient creates a new SpanAPI client based on the command line options and/or
// defaults.
func NewSpanAPIClient() (*spanapi.APIClient, context.Context, context.CancelFunc) {
	config := spanapi.NewConfiguration()
	config.Debug = global.Options.Debug
	if global.Options.OverrideEndpoint != "" {
		config.Servers = spanapi.ServerConfigurations{
			spanapi.ServerConfiguration{URL: global.Options.OverrideEndpoint},
		}
	}
	ctx, done := apitools.ContextWithAuthAndTimeout(global.Options.Token, global.Options.Timeout)
	client := spanapi.NewAPIClient(config)
	CheckVersion(ctx, client)

	return client, ctx, done
}
