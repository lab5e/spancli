package params

import (
	"github.com/jessevdk/go-flags"
	"github.com/lab5e/spancli/pkg/cmd/cert"
	"github.com/lab5e/spancli/pkg/cmd/collection"
	"github.com/lab5e/spancli/pkg/cmd/device"
	"github.com/lab5e/spancli/pkg/cmd/firmware"
	"github.com/lab5e/spancli/pkg/cmd/inbox"
	"github.com/lab5e/spancli/pkg/cmd/invite"
	"github.com/lab5e/spancli/pkg/cmd/outbox"
	"github.com/lab5e/spancli/pkg/cmd/output"
	"github.com/lab5e/spancli/pkg/cmd/team"
	"github.com/lab5e/spancli/pkg/cmd/version"
	"github.com/lab5e/spancli/pkg/global"
)

// A bit of a redirect here to avoid cyclic dependencies; the commands need to know the global parameters
// so a separate variable in another package holds the globals.

type Options struct {
	Args *global.Parameters

	Team       team.Command       `command:"team" description:"team management"`
	Invite     invite.Command     `command:"invite" description:"manage team invitations"`
	Collection collection.Command `command:"collection" alias:"col" description:"collection management"`
	Device     device.Command     `command:"device" alias:"dev" description:"device management"`
	Inbox      inbox.Command      `command:"inbox" description:"Read messages from devices"`
	Outbox     outbox.Command     `command:"outbox" description:"Send messages to devices"`
	Output     output.Command     `command:"output" alias:"out" description:"output management"`
	Cert       cert.Command       `command:"cert" description:"certificate management"`
	Firmware   firmware.Command   `command:"fw" description:"firmware management"`
	Version    version.Command    `command:"version" description:"show version"`
}

var opt Options
var Parser *flags.Parser

func init() {
	opt.Args = global.Options
	Parser = flags.NewParser(&opt, flags.Default)
}