package login

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lab5e/go-spanuserapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

// Command is the subcommand to log in. Logging in saves a JWT token locally (in ~/.spanrc) which will be
// used to authenticate with Span.
type Command struct {
	Username string `long:"username" short:"u" description:"user name"`
	Password string `long:"password" short:"p" description:"password"`
	Passcode string `long:"passcode" short:"c" description:"MFA token passcode"`
}

func (c *Command) Execute([]string) error {
	userApi := helpers.NewUserAPIClient()
	ctx, done := context.WithTimeout(context.Background(), time.Second*60)
	defer done()
	res, httpRes, err := userApi.SessionApi.Login(ctx).Body(spanuserapi.LoginRequest{
		Email:    &c.Username,
		Password: &c.Password,
		Passcode: &c.Passcode,
	}).Execute()

	if err != nil {
		return helpers.APIError(httpRes, err)
	}

	if *res.Result == spanuserapi.LOGINRESPONSERESULT_INCOMPLETE {
		fmt.Println("Needs MFA passcode to log in")
		return errors.New("missing MFA passcode")
	}
	if err := helpers.WriteCredentials(*res.Credentials); err != nil {
		return err
	}
	fmt.Println("Logged in")
	return nil
}
