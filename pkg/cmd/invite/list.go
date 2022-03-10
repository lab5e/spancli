package invite

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listInvite struct {
	TeamID   string `long:"team-id" description:"id of team we wish to list invites for" required:"yes"`
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

func (r *listInvite) Execute([]string) error {
	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	invites, res, err := client.TeamsApi.ListInvites(ctx, r.TeamID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	if r.Format == "json" {
		json, err := json.MarshalIndent(invites, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	// place it after the JSON output since outputting an empty JSON object
	// is better than outputting a string.
	if invites.Invites == nil {
		fmt.Printf("no invites\n")
		return nil
	}

	t := helpers.NewTableOutput(r.Format, r.NoColor, r.PageSize)
	t.SetTitle("Invites for team " + r.TeamID)
	t.AppendHeader(table.Row{"Code", "Created"})
	for _, invite := range *invites.Invites {
		createdAt := helpers.LocalTimeFormat(*invite.CreatedAt)
		t.AppendRow(table.Row{*invite.Code, createdAt})
	}
	helpers.RenderTable(t, r.Format)

	return nil
}
