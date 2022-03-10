package outbox

import (
	"fmt"
	"strings"
)

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
