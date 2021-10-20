package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jessevdk/go-flags"
	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/go-spanapi/v4/apitools"
	"github.com/lab5e/go-userapi"
	userapitools "github.com/lab5e/go-userapi/apitools"
	"github.com/lab5e/spancli/pkg/helpers"
)

type options struct {
	Token   string        `long:"token" env:"SPAN_API_TOKEN" description:"span API token" required:"yes"`
	Timeout time.Duration `long:"timeout" default:"15s" description:"timeout for operation"`
	Debug   bool          `long:"debug" description:"turn on debug output"`

	Team       teamCmd       `command:"team" description:"team management"`
	Invite     inviteCmd     `command:"invite" description:"manage team invitations"`
	Collection collectionCmd `command:"collection" alias:"col" description:"collection management"`
	Device     deviceCmd     `command:"device" alias:"dev" description:"device management"`
	Data       dataCmd       `command:"data" description:"data listing commands"`
	Version    versionCmd    `command:"version" description:"show version"`
}

var opt options
var parser = flags.NewParser(&opt, flags.Default)

const timeFmt = "2006-01-02 15:04:05"

func main() {
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			switch flagsErr.Type {
			case flags.ErrHelp:
				os.Exit(0)

			case flags.ErrCommandRequired:
				os.Exit(1)

			case flags.ErrUnknownCommand:
				os.Exit(1)

			case flags.ErrRequired:
				os.Exit(1)

			case flags.ErrUnknownFlag:
				os.Exit(1)

			case flags.ErrMarshal:
				os.Exit(1)

			case flags.ErrExpectedArgument:
				os.Exit(1)

			default:
				fmt.Printf("%v [%d]\n", err, flagsErr.Type)
				os.Exit(0)
			}
		}
		os.Exit(1)
	}
}

// apiError creates an error instance based on error message
// and HTTP response returned from API call.
func apiError(res *http.Response, e error) error {
	var body []byte
	var err error

	if res != nil {
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			// If we can't read the body, just return the error from the API
			return e
		}
	}

	var errmsg struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	err = json.Unmarshal(body, &errmsg)
	if err != nil {
		// If we can't extract the error message we just return the error from the API
		return e
	}

	return fmt.Errorf("%s: %s", res.Status, errmsg.Message)
}

// newSpanAPIClient creates a new SpanAPI client based on the command line options and/or
// defaults.
func newSpanAPIClient() (*spanapi.APIClient, context.Context, context.CancelFunc) {
	config := spanapi.NewConfiguration()
	config.Debug = opt.Debug

	ctx, done := apitools.ContextWithAuthAndTimeout(opt.Token, opt.Timeout)
	client := spanapi.NewAPIClient(config)
	helpers.CheckVersion(ctx, client)

	return client, ctx, done
}

// newUserAPIClient creates a new UserAPI client based on the command line options
// and/or defaults.
func newUserAPIClient() (*userapi.APIClient, context.Context, context.CancelFunc) {
	config := userapi.NewConfiguration()
	config.Debug = opt.Debug

	ctx, done := userapitools.NewAuthenticatedContext(opt.Token, opt.Timeout)
	client := userapi.NewAPIClient(config)

	return client, ctx, done
}

// verifyDeleteIntent requires the user to type in "yes " followed by
// a random number so as to avoid accidental deletion.  This can be
// overridden by including the --yes-i-am-sure flag on the command
// line.
func verifyDeleteIntent() bool {
	rand.Seed(time.Now().UnixNano())
	verify := fmt.Sprintf("yes %04d", rand.Intn(9999))

	fmt.Printf("\n%s\n\n", text.Colors{text.BgRed, text.FgWhite}.Sprint("*** D A N G E R ***"))
	fmt.Printf("enter '%s' to confirm: ", verify)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	return text == (verify + "\n")
}

func truncateString(s string, n int) string {
	if len(s) > n && len(s) > 3 {
		return s[:n-3] + "..."
	}
	return s
}

func msSinceEpochToTime(ts string) (int64, time.Time) {
	r, err := strconv.ParseInt(ts, 10, 63)
	if err != nil {
		return time.Now().UnixNano() / int64(time.Millisecond), time.Now()
	}
	return r, time.Unix(0, r*int64(time.Millisecond))
}

func localTimeFormat(ts string) string {
	_, t := msSinceEpochToTime(ts)
	return t.Local().Format(timeFmt)
}

func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
