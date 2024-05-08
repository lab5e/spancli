package logout

import (
	"context"
	"fmt"
	"time"

	"github.com/lab5e/spancli/pkg/helpers"
)

// Command is the subcommand to log in. Logging in saves a JWT token locally (in ~/.spanrc) which will be
// used to authenticate with Span.
type Command struct {
}

func (c *Command) Execute([]string) error {
	creds := helpers.ReadCredentials()
	if creds == "" {
		helpers.RemoveCredentials()
		return nil
	}
	userApi := helpers.NewUserAPIClient()
	ctx, done := context.WithTimeout(context.Background(), time.Second*60)
	defer done()
	_, httpRes, err := userApi.SessionApi.Logout(ctx).Body(make(map[string]interface{})).Execute()

	// Ignore if status code is < 500; the token might be invalid
	if err != nil && httpRes.StatusCode > 499 {
		helpers.APIError(httpRes, err)
		return err
	}
	helpers.RemoveCredentials()
	fmt.Println("Logged out")
	return nil
}
