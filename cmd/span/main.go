package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/lab5e/spanclient-go/v4"
)

type options struct {
	Token   string `long:"token" env:"SPAN_API_TOKEN" description:"span API token" required:"yes"`
	Timeout int    `long:"timeout" default:"120" description:"timeout in number of seconds"`
	Debug   bool   `long:"debug" description:"turn on debug output"`

	Collection collectionCmd `command:"collection" alias:"col" description:"collection management"`
	Device     deviceCmd     `command:"device" alias:"dev" description:"device management"`
	Data       dataCmd       `command:"data" description:"data listing commands"`
	Listen     listenCmd     `command:"listen" description:"live streaming of data"`
}

var opt options
var parser = flags.NewParser(&opt, flags.Default)

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

			default:
				fmt.Printf("%v [%d]\n", err, flagsErr.Type)
				os.Exit(0)
			}
		}
		// The error will be printed by the flags library so no need to
		// print it again
		os.Exit(1)
	}
}

// Create spanclient.Configuration based on the command line options.
func clientConfig() *spanclient.Configuration {
	config := spanclient.NewConfiguration()
	config.Debug = opt.Debug

	// For debugging purposes. This *could* be an option on the command
	// itself but nobody but one or two people in the world (including me)
	// needs this feature once in a while so we'll just keep it as an
	// environment variable.
	apiAddr := os.Getenv("SPAN_HOST")
	if apiAddr != "" {
		config.BasePath = apiAddr
	}
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

func truncateString(s string, n int) string {
	if len(s) > n && len(s) > 3 {
		return s[:n-3] + "..."
	}
	return s
}

// tagsToString converts a tag struct to a name:value list of strings. The name
// tag (if it is set) is placed first since it makes it easier to read the list
func tagsToString(tags map[string]string) string {
	if tags == nil {
		return ""
	}
	var ret []string
	name, hasName := tags["name"]
	if hasName {
		ret = append(ret, fmt.Sprintf("name:%s", name))
	}

	for name, value := range tags {
		if name == "name" {
			continue
		}
		ret = append(ret, fmt.Sprintf("%s:%s", name, value))
	}

	return strings.Join(ret, "  ")
}
