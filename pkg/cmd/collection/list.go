package collection

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listCollection struct {
	Format commonopt.ListFormat
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

	if r.Format.Format == "json" {
		json, err := json.MarshalIndent(collections, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := helpers.NewTableOutput(r.Format)
	t.SetTitle("Collections")
	t.AppendHeader(table.Row{
		"ID",
		"Name",
		"Firmware",
		"Target Image",
		"Tags",
		"TeamID",
	})
	for _, col := range collections.Collections {
		// only truncate name if we output as 'text'
		name := col.GetTags()["name"]
		if r.Format.Format == "text" {
			name = helpers.EllipsisString(name, 25)
		}
		t.AppendRow(table.Row{
			col.GetCollectionId(),
			name,
			col.Firmware.GetManagement(),
			col.Firmware.GetTargetFirmwareId(),
			helpers.TagsToString(col.GetTags()),
			col.GetTeamId(),
		})
	}
	helpers.RenderTable(t, r.Format.Format)
	return nil
}
