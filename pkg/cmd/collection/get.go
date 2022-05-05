package collection

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type getCollection struct {
	ID     commonopt.Collection
	Format commonopt.ListFormat
}

func (r *getCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	col, res, err := client.CollectionsApi.RetrieveCollection(ctx, r.ID.CollectionID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	t := helpers.NewTableOutput(r.Format)

	t.SetTitle("Collection %s", r.ID.CollectionID)

	t.AppendHeader(table.Row{"Field", "Value"})

	if err := helpers.DumpToTable(t, col); err != nil {
		return err
	}

	helpers.RenderTable(t, r.Format.Format)
	return nil
}
