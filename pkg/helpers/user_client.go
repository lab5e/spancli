package helpers

import (
	"context"

	"github.com/lab5e/go-userapi"
	userapitools "github.com/lab5e/go-userapi/apitools"
	"github.com/lab5e/spancli/pkg/global"
)

// newUserAPIClient creates a new UserAPI client based on the command line options
// and/or defaults.
func NewUserAPIClient() (*userapi.APIClient, context.Context, context.CancelFunc) {
	config := userapi.NewConfiguration()
	config.Debug = global.Options.Debug

	ctx, done := userapitools.NewAuthenticatedContext(global.Options.Token, global.Options.Timeout)
	client := userapi.NewAPIClient(config)

	return client, ctx, done
}
