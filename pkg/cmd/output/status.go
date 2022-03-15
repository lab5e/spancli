package output

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type outputStatus struct {
	ID     commonopt.Collection
	OID    oid
	Format commonopt.ListFormat
}

func (c *outputStatus) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	status, res, err := client.OutputsApi.Status(ctx, c.ID.CollectionID, c.OID.OutputID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	t := helpers.NewTableOutput(c.Format)
	t.AppendHeader(table.Row{
		"Enabled",
		"Error count",
		"Forwarded",
		"Received",
		"Retransmits",
	})

	t.AppendRow(table.Row{
		status.GetEnabled(),
		status.GetErrorCount(),
		status.GetForwarded(),
		status.GetReceived(),
		status.GetRetransmits(),
	})
	helpers.RenderTable(t, c.Format.Format)
	return nil
}
