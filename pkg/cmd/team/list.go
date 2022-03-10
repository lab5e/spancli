package team

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listTeams struct {
	//lint:ignore SA5008 Linter is unhappy with multiple choice values
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

func (r *listTeams) Execute([]string) error {
	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	teamList, res, err := client.TeamsApi.ListTeams(ctx).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	if r.Format == "json" {
		json, err := json.MarshalIndent(*teamList.Teams, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := helpers.NewTableOutput(r.Format, r.NoColor, r.PageSize)
	t.SetTitle("Teams")
	t.AppendHeader(table.Row{"ID", "Name"})

	for _, team := range *teamList.Teams {
		// only truncate name if we output as 'text'
		name := team.GetTags()["name"]
		if r.Format == "text" {
			name = helpers.EllipsisString(name, 25)
		}

		if team.GetIsPrivate() {
			name += " [P]"
		}

		t.AppendRow(table.Row{*team.TeamId, name})
	}
	helpers.RenderTable(t, r.Format)

	return nil
}
