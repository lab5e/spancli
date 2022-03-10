package outbox

import "fmt"

type watchOutboxCmd struct {
	MessageID string `long:"id" description:"message ID for message to watch; empty id watches entire outbox"`
	Format    string `long:"format" default:"text" description:"which output format to use" choice:"text" choice:"json"`
}

func (*watchOutboxCmd) Execute([]string) error {
	fmt.Println("Watch outbox")
	return nil
}
