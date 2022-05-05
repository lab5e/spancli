package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/go-spanapi/v4/apitools"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/global"
	"github.com/lab5e/spancli/pkg/helpers"
)

// Watch outbox for changes. This is just a simple polling operation every 2 seconds to check for changes in the mailbox
type watchOutboxCmd struct {
	ID        commonopt.CollectionAndDevice
	MessageID string `long:"id" description:"message ID for message to watch; empty id watches entire outbox"`
}

func (c *watchOutboxCmd) Execute([]string) error {
	client, _, ctxDone := helpers.NewSpanAPIClient()
	ctxDone()

	fmt.Printf("Polling outbox for collection %s/device %s...\n", c.ID.CollectionID, c.ID.DeviceID)

	authContext := apitools.ContextWithAuth(global.Options.Token)

	outboxState := make(map[string]spanapi.MessageState)
	populate := true
	for {
		ctx, done := context.WithTimeout(authContext, 5*time.Second)
		list, res, err := client.DevicesApi.ListDownstreamMessages(ctx, c.ID.CollectionID, c.ID.DeviceID).Execute()
		if err != nil {
			fmt.Printf("Got error polling outbox: %v. Stopping", err)
			done()
			return helpers.APIError(res, err)
		}
		if populate {
			for _, msg := range list.Messages {
				outboxState[msg.GetMessageId()] = msg.GetState()
			}
			populate = false
		}
		for _, msg := range list.Messages {
			id := msg.GetMessageId()
			state := msg.GetState()

			existing, ok := outboxState[id]
			if !ok {
				// New event
				fmt.Printf("%s is added to the outbox\n", id)
			} else {
				if existing != state {
					fmt.Printf("%s changed state from %s to %s\n", id, existing, state)
				}
			}
			// Update state
			outboxState[id] = state

		}
		done()
		time.Sleep(2 * time.Second)
	}
}
