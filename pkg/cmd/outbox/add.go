package outbox

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type addOutboxCmd struct {
	ID commonopt.CollectionAndDevice
	//lint:ignore SA5008 Multiple choices makes the linter unhappy
	Format string `long:"format" description:"input format" default:"base64" choice:"text" choice:"base64" choice:"hex"`
}

func (o *addOutboxCmd) Execute(params []string) error {
	rawPayload := strings.TrimSpace(strings.Join(params, " "))
	if rawPayload == "" {
		return errors.New("must specify payload")
	}
	var b64Payload string
	switch o.Format {
	case "text":
		b64Payload = base64.StdEncoding.EncodeToString([]byte(rawPayload))
	case "base64":
		b64Payload = rawPayload
	case "hex":
		bytes, err := hex.DecodeString(rawPayload)
		if err != nil {
			return err
		}
		b64Payload = base64.StdEncoding.EncodeToString(bytes)
	}

	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	msg, res, err := client.DevicesApi.AddDownstreamMessage(ctx, o.ID.CollectionID, o.ID.DeviceID).Body(
		spanapi.AddDownstreamMessageRequest{
			Payload: spanapi.PtrString(b64Payload),
		},
	).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("Message with ID %s added to outbox for device %s in collection %s\n", msg.GetMessageId(), msg.GetDeviceId(), msg.GetCollectionId())
	return nil
}
