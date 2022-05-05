package outbox

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listOutboxCmd struct {
	ID     commonopt.CollectionAndDevice
	List   commonopt.QueryOptions
	Format commonopt.ListFormat
}

func (c *listOutboxCmd) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	req := client.DevicesApi.ListDownstreamMessages(ctx, c.ID.CollectionID, c.ID.DeviceID)
	if c.List.Start != "" {
		req = req.Start(c.List.Start)
	}
	if c.List.End != "" {
		req = req.End(c.List.End)
	}
	if c.List.Limit != 0 {
		req = req.Limit(c.List.Limit)
	}

	list, res, err := req.Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}
	t := helpers.NewTableOutput(c.Format)
	t.SetTitle(fmt.Sprintf("Inbox for device %s on collection %s", c.ID.DeviceID, c.ID.CollectionID))
	t.AppendHeader(table.Row{"Message ID", "Created", "State", "Sent", "Transport", "Payload"})
	for _, item := range list.GetMessages() {
		t.AppendRow(table.Row{
			item.GetMessageId(),
			helpers.DateFormat(item.GetCreatedTime(), c.Format.NumericDate),
			item.GetState(),
			helpers.DateFormat(item.GetSentTime(), c.Format.NumericDate),
			item.GetTransport(),
			helpers.EllipsisString(helpers.PayloadFormat(item.GetPayload(), c.List.Decode), c.Format.MaxPayloadWdith),
		})
	}
	helpers.RenderTable(t, c.Format.Format)
	return nil
}
