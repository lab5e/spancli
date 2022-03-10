package team

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listMembers struct {
	TeamID   string `long:"team-id" description:"id of team" required:"yes"`
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

func (r *listMembers) Execute([]string) error {
	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	team, res, err := client.TeamsApi.RetrieveTeam(ctx, r.TeamID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	if r.Format == "json" {
		json, err := json.MarshalIndent(team.Members, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := helpers.NewTableOutput(r.Format, r.NoColor, r.PageSize)
	t.SetTitle("Members of " + r.TeamID)
	t.AppendHeader(table.Row{"UserID", "Role", "Email"})

	for _, member := range *team.Members {
		t.AppendRow(table.Row{*member.UserId, *member.Role, *member.User.Email})
	}
	helpers.RenderTable(t, r.Format)
	return nil
}
