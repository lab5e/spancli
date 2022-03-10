package collection

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listCollection struct {
	//lint:ignore SA5008 Multiple choice tags makes the linter unhappy
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

func (r *listCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	collections, res, err := client.CollectionsApi.ListCollections(ctx).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	if collections.Collections == nil {
		fmt.Printf("no collections\n")
		return nil
	}

	if r.Format == "json" {
		json, err := json.MarshalIndent(collections, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := helpers.NewTableOutput(r.Format, r.NoColor, r.PageSize)
	t.SetTitle("Collections")
	t.AppendHeader(table.Row{"ID", "Name", "TeamID"})
	for _, col := range collections.Collections {
		// only truncate name if we output as 'text'
		name := col.GetTags()["name"]
		if r.Format == "text" {
			name = helpers.EllipsisString(name, 25)
		}

		t.AppendRow(table.Row{
			*col.CollectionId,
			name,
			*col.TeamId,
		})
	}
	helpers.RenderTable(t, r.Format)
	return nil
}
