package output

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type outputLogs struct {
	ID     commonopt.Collection
	OID    oid
	Format commonopt.ListFormat
}

func (c *outputLogs) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	list, res, err := client.OutputsApi.Logs(ctx, c.ID.CollectionID, c.OID.OutputID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	t := helpers.NewTableOutput(c.Format)

	t.SetTitle("Logs for output %s", c.OID.OutputID)
	t.AppendHeader(table.Row{
		"Time",
		"Message",
		"Repeats",
	})

	for _, l := range list.Logs {
		t.AppendRow(table.Row{
			l.GetTime(),
			l.GetMessage(),
			l.GetRepeated(),
		})
	}
	helpers.RenderTable(t, c.Format.Format)
	return nil
}
