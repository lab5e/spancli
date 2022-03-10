package outbox

import (
	"fmt"
	"strings"
)

type Command struct {
	CollectionID string          `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string          `long:"device-id" description:"device id"`
	Add          addOutboxCmd    `command:"add" description:"add message to oubox"`
	List         listOutboxCmd   `command:"list" description:"list messages in outbox"`
	Delete       deleteOutboxCmd `command:"delete" description:"delete message from outbox"`
	Watch        watchOutboxCmd  `command:"watch" description:"monitor message in outbox"`
}
type addOutboxCmd struct {
	Text   bool `long:"text" description:"Payload, plain text format"`
	Base64 bool `long:"base64" description:"Payload, base64 encoded"`
	Hex    bool `long:"hex" description:"Payload, hex encoded"`
}

func (*addOutboxCmd) Execute(params []string) error {
	payload := strings.Join(params, " ")
	fmt.Println("add outbox ", payload)
	return nil
}

type listOutboxCmd struct {
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

func (*listOutboxCmd) Execute([]string) error {
	fmt.Println("list outbox")
	return nil
}

type deleteOutboxCmd struct {
	MessageID string `long:"id" description:"message ID to delete"`
}

func (*deleteOutboxCmd) Execute([]string) error {
	fmt.Println("Delete in outbox")
	return nil
}

type watchOutboxCmd struct {
	MessageID string `long:"id" description:"message ID for message to watch; empty id watches entire outbox"`
	Format    string `long:"format" default:"text" description:"which output format to use" choice:"text" choice:"json"`
}

func (*watchOutboxCmd) Execute([]string) error {
	fmt.Println("Watch outbox")
	return nil
}
