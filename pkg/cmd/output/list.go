package output

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listOutput struct {
	ID     commonopt.Collection
	Format commonopt.ListFormat
}

func (c *listOutput) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	list, res, err := client.OutputsApi.ListOutputs(ctx, c.ID.CollectionID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	t := helpers.NewTableOutput(c.Format)
	t.SetTitle("Outputs for collection %s", c.ID.CollectionID)
	t.AppendHeader(table.Row{
		"Output ID",
		"Type",
		"Enabled",
		"Config",
		"Tags",
	})
	for _, o := range list.Outputs {
		t.AppendRow(table.Row{
			o.GetOutputId(),
			o.GetType(),
			o.GetEnabled(),
			configToString(o.Config),
			helpers.TagsToString(o.GetTags()),
		})
	}
	helpers.RenderTable(t, c.Format.Format)
	return nil
}

// Just convert the config to a single-line JSON string
func configToString(cfg *spanapi.OutputConfig) string {
	if cfg == nil {
		return ""
	}
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		fmt.Printf("\nError marshaling config: %v\n", err)
		return ""
	}
	return string(jsonStr)
}
