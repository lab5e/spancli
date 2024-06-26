package params

import (
	"github.com/jessevdk/go-flags"
	"github.com/lab5e/spancli/pkg/cmd/activity"
	"github.com/lab5e/spancli/pkg/cmd/blob"
	"github.com/lab5e/spancli/pkg/cmd/cert"
	"github.com/lab5e/spancli/pkg/cmd/collection"
	"github.com/lab5e/spancli/pkg/cmd/device"
	"github.com/lab5e/spancli/pkg/cmd/firmware"
	"github.com/lab5e/spancli/pkg/cmd/gateway"
	"github.com/lab5e/spancli/pkg/cmd/inbox"
	"github.com/lab5e/spancli/pkg/cmd/login"
	"github.com/lab5e/spancli/pkg/cmd/logout"
	"github.com/lab5e/spancli/pkg/cmd/outbox"
	"github.com/lab5e/spancli/pkg/cmd/output"
	"github.com/lab5e/spancli/pkg/cmd/sample"
	"github.com/lab5e/spancli/pkg/cmd/version"
	"github.com/lab5e/spancli/pkg/global"
)

// A bit of a redirect here to avoid cyclic dependencies; the commands need to know the global parameters
// so a separate variable in another package holds the globals.

// Options is the main struct for the span command. Each command is put into this.
type Options struct {
	Args *global.Parameters

	Collection collection.Command `command:"collection" alias:"col" description:"collection management"`
	Device     device.Command     `command:"device" alias:"dev" description:"device management"`
	Inbox      inbox.Command      `command:"inbox" description:"Read messages from devices"`
	Outbox     outbox.Command     `command:"outbox" description:"Send messages to devices"`
	Output     output.Command     `command:"output" alias:"out" description:"output management"`
	Cert       cert.Command       `command:"cert" description:"certificate management"`
	Firmware   firmware.Command   `command:"fw" description:"firmware management"`
	Blob       blob.Command       `command:"blob" description:"blob maanagement"`
	Login      login.Command      `command:"login" description:"Autenticate with the service"`
	Logout     logout.Command     `command:"logout" description:"Log out from the service"`
	Samples    sample.Command     `command:"samples" alias:"ex" description:"Sample code"`
	Version    version.Command    `command:"version" description:"show version"`
	Activity   activity.Command   `command:"activity" description:"monitor activity"`
	Gateway    gateway.Command    `command:"gateway" alias:"gw" description:"gateway management"`
}

var opt Options

// Parser is the comand line parser used to parse the options
var Parser *flags.Parser

func init() {
	opt.Args = global.Options
	Parser = flags.NewParser(&opt, flags.Default)
}
