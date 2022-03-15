package team

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listTeams struct {
	Format commonopt.ListFormat
}

func (r *listTeams) Execute([]string) error {
	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	teamList, res, err := client.TeamsApi.ListTeams(ctx).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	if r.Format.Format == "json" {
		json, err := json.MarshalIndent(*teamList.Teams, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := helpers.NewTableOutput(r.Format)
	t.SetTitle("Teams")
	t.AppendHeader(table.Row{"ID", "Name", "Tags"})

	for _, team := range *teamList.Teams {
		// only truncate name if we output as 'text'
		name := team.GetTags()["name"]
		if r.Format.Format == "text" {
			name = helpers.EllipsisString(name, 25)
		}

		if team.GetIsPrivate() {
			name += " [P]"
		}

		t.AppendRow(table.Row{
			team.GetTeamId(),
			name,
			helpers.TagsToString(team.GetTags()),
		})
	}
	helpers.RenderTable(t, r.Format.Format)

	return nil
}
