package team

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listMembers struct {
	TeamID string `long:"team-id" description:"id of team" required:"yes"`

	Format commonopt.ListFormat
}

func (r *listMembers) Execute([]string) error {
	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	team, res, err := client.TeamsApi.RetrieveTeam(ctx, r.TeamID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	if r.Format.Format == "json" {
		json, err := json.MarshalIndent(team.Members, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := helpers.NewTableOutput(r.Format)
	t.SetTitle("Members of " + r.TeamID)
	t.AppendHeader(table.Row{"UserID", "Role", "Email"})

	for _, member := range *team.Members {
		t.AppendRow(table.Row{*member.UserId, *member.Role, *member.User.Email})
	}
	helpers.RenderTable(t, r.Format.Format)
	return nil
}
