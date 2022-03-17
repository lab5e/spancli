package output

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type getOutput struct {
	ID     commonopt.Collection
	OID    oid
	Format commonopt.ListFormat
}

func (c *getOutput) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	fw, res, err := client.OutputsApi.RetrieveOutput(ctx, c.ID.CollectionID, c.OID.OutputID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	t := helpers.NewTableOutput(c.Format)

	t.SetTitle("Output %s", c.OID.OutputID)

	t.AppendHeader(table.Row{"Field", "Value"})

	if err := helpers.DumpToTable(t, fw); err != nil {
		return err
	}

	helpers.RenderTable(t, c.Format.Format)
	return nil
}
