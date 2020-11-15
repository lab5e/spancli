package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/lab5e/spanclient-go/v4"
)

type options struct {
	Token   string `long:"token" env:"SPAN_API_TOKEN" description:"span API token" required:"yes"`
	Timeout int    `long:"timeout" default:"120" description:"timeout in number of seconds"`
	Debug   bool   `long:"debug" description:"turn on debug output"`
}

var opt options
var parser = flags.NewParser(&opt, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

// Create spanclient.Configuration based on the command line options.
func clientConfig() *spanclient.Configuration {
	config := spanclient.NewConfiguration()
	config.Debug = opt.Debug

	return config
}

// spanContextWithTimeout creates a context.Context with timeout and
// credentials
func spanContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithValue(context.Background(),
		spanclient.ContextAPIKey,
		spanclient.APIKey{
			Key:    opt.Token,
			Prefix: "",
		}),
		time.Duration(opt.Timeout)*time.Second)
}

// verifyDeleteIntent requires the user to type in "yes " followed by
// a random number so as to avoid accidental deletion.  This can be
// overridden by including the --yes-i-am-sure flag on the command
// line.
func verifyDeleteIntent() bool {
	rand.Seed(time.Now().UnixNano())
	verify := fmt.Sprintf("yes %04d", rand.Intn(9999))

	fmt.Printf("\n*** D A N G E R ***\n\nenter '%s' to confirm: ", verify)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	return text == (verify + "\n")
}

// timeToMilliseconds converts a time.Time to milliseconds since epoch
func timeToMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func truncateString(s string, n int) string {
	if len(s) > n && len(s) > 3 {
		return s[:n-3] + "..."
	}
	return s
}
