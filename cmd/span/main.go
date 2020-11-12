package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/lab5e/spanclient-go"
)

type options struct {
	Token   string `long:"token" env:"SPAN_API_TOKEN" description:"span API token" required:"yes"`
	Timeout int    `long:"timeout" description:"timeout in number of seconds"`
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
	config.Scheme = "https"

	log.Printf("Config: %+v", config)

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
