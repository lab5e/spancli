package login

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/lab5e/go-spanuserapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
	"golang.org/x/term"
)

// Command is the subcommand to log in. Logging in saves a JWT token locally (in ~/.spanrc) which will be
// used to authenticate with Span.
type Command struct {
	Email    string `long:"email" short:"u" description:"email address name"`
	Password string `long:"password" short:"p" description:"password"`
	Passcode string `long:"passcode" short:"c" description:"MFA token passcode"`
}

func (c *Command) Execute([]string) error {
	if c.Email == "" {
		c.Email = prompt("Email: ")
	}

	if c.Password == "" {
		c.Password = prompt("Password: ", true)
	}

	userApi := helpers.NewUserAPIClient()
	ctx, done := context.WithTimeout(context.Background(), time.Second*60)
	defer done()
	res, httpRes, err := userApi.SessionApi.Login(ctx).Body(spanuserapi.LoginRequest{
		Email:    &c.Email,
		Password: &c.Password,
		Passcode: &c.Passcode,
	}).Execute()

	if err != nil {
		return helpers.APIError(httpRes, err)
	}

	if *res.Result == spanuserapi.LOGINRESPONSERESULT_INCOMPLETE {
		fmt.Println("Needs MFA passcode to log in, please add -c <passcode> or --passcode <passcode>")
		return errors.New("missing MFA passcode")
	}
	if err := helpers.WriteCredentials(*res.Credentials); err != nil {
		return err
	}

	fmt.Printf("\n\nLogged in as %s\n\n", c.Email)

	return nil
}

func prompt(s string, password ...bool) string {
	fmt.Print(s)
	if len(password) > 0 && password[0] {
		inputBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(inputBytes))
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	return strings.TrimSpace(input)
}
