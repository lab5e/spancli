package outbox

import "fmt"

type deleteOutboxCmd struct {
	MessageID string `long:"id" description:"message ID to delete"`
}

func (*deleteOutboxCmd) Execute([]string) error {
	fmt.Println("Delete in outbox")
	return nil
}
