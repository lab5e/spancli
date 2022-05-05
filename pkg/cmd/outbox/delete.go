package outbox

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteOutboxCmd struct {
	ID        commonopt.CollectionAndDevice
	MessageID string `long:"id" description:"message ID to delete"`
}

func (c *deleteOutboxCmd) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	msg, res, err := client.DevicesApi.DeleteDownstreamMessage(ctx, c.ID.CollectionID, c.ID.DeviceID, c.MessageID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("Removed message with ID %s from outbox\n", msg.GetMessageId())
	return nil
}
