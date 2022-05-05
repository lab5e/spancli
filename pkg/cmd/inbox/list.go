package inbox

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listInboxCmd struct {
	ID     commonopt.CollectionAndOptionalDevice
	List   commonopt.QueryOptions
	Format commonopt.ListFormat
}

func (c *listInboxCmd) Execute([]string) error {
	if c.ID.DeviceID == "" {
		return c.listCollectionData()
	}
	return c.listDeviceInbox()
}

func (c *listInboxCmd) listDeviceInbox() error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	req := client.DevicesApi.ListUpstreamMessages(ctx, c.ID.CollectionID, c.ID.DeviceID)
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
	t.AppendHeader(table.Row{"Message ID", "Received", "Transport", "Payload"})
	for _, item := range list.GetMessages() {
		t.AppendRow(table.Row{
			item.GetMessageId(),
			helpers.DateFormat(item.GetReceived(), c.Format.NumericDate),
			item.GetTransport(),
			helpers.EllipsisString(helpers.PayloadFormat(item.GetPayload(), c.List.Decode), c.Format.MaxPayloadWdith),
		})
	}
	helpers.RenderTable(t, c.Format.Format)
	return nil
}

func (c *listInboxCmd) listCollectionData() error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	req := client.CollectionsApi.ListCollectionData(ctx, c.ID.CollectionID)
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
	t.SetTitle(fmt.Sprintf("Inbox for collection %s", c.ID.CollectionID))
	t.AppendHeader(table.Row{"Message ID", "Device ID", "Received", "Transport", "Payload"})
	for _, item := range list.GetData() {
		t.AppendRow(table.Row{
			item.GetMessageId(),
			*item.GetDevice().DeviceId,
			helpers.DateFormat(item.GetReceived(), c.Format.NumericDate),
			item.GetTransport(),
			helpers.EllipsisString(helpers.PayloadFormat(item.GetPayload(), c.List.Decode), c.Format.MaxPayloadWdith),
		})
	}
	helpers.RenderTable(t, c.Format.Format)
	return nil
}
