package inbox

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listInboxCmd struct {
	ID     commonopt.CollectionAndOptionalDevice
	List   commonopt.ListOptions
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
		return helpers.ApiError(res, err)
	}
	t := helpers.NewTableOutput(c.Format)
	t.SetTitle(fmt.Sprintf("Inbox for device %s on collection %s", c.ID.DeviceID, c.ID.CollectionID))
	t.AppendHeader(table.Row{"Message ID", "Received", "Transport", "Payload"})
	for _, item := range list.GetMessages() {
		t.AppendRow(table.Row{
			item.GetMessageId(),
			dateFormat(item.GetReceived(), c.Format.NumericDate),
			item.GetTransport(),
			helpers.EllipsisString(payloadFormat(item.GetPayload(), c.List.Decode), c.Format.MaxPayloadWdith),
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
		return helpers.ApiError(res, err)
	}
	t := helpers.NewTableOutput(c.Format)
	t.SetTitle(fmt.Sprintf("Inbox for collection %s", c.ID.CollectionID))
	t.AppendHeader(table.Row{"Message ID", "Device ID", "Received", "Transport", "Payload"})
	for _, item := range list.GetData() {
		t.AppendRow(table.Row{
			item.GetMessageId(),
			*item.GetDevice().DeviceId,
			dateFormat(item.GetReceived(), c.Format.NumericDate),
			item.GetTransport(),
			helpers.EllipsisString(payloadFormat(item.GetPayload(), c.List.Decode), c.Format.MaxPayloadWdith),
		})
	}
	helpers.RenderTable(t, c.Format.Format)
	return nil
}

func dateFormat(dateStr string, numeric bool) string {
	if numeric {
		return dateStr
	}

	val, err := strconv.ParseInt(dateStr, 10, 64)
	if err == nil {
		return time.UnixMilli(val).Format(time.RFC3339)
	}
	return "(error)"
}
func payloadFormat(pl string, decode bool) string {
	if decode {
		ret, err := base64.StdEncoding.DecodeString(pl)
		if err != nil {
			return "(error)"
		}
		return string(ret)
	}
	return pl
}
