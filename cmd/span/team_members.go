package main

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type memberCmd struct {
	List   listMembers  `command:"list" alias:"ls" description:""`
	Delete deleteMember `command:"delete" alias:"del" description:""`
}

type listMembers struct {
	TeamID   string `long:"team-id" description:"id of team" required:"yes"`
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

type deleteMember struct {
	TeamID     string `long:"team-id" description:"id of team" required:"yes"`
	UserID     string `long:"user-id" description:"id of user we want to remove" required:"yes"`
	YesIAmSure bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

func (r *listMembers) Execute([]string) error {
	client, ctx, cancel := newUserAPIClient()
	defer cancel()

	team, res, err := client.TeamsApi.RetrieveTeam(ctx, r.TeamID).Execute()
	if err != nil {
		return apiError(res, err)
	}

	if r.Format == "json" {
		json, err := json.MarshalIndent(team.Members, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := newTableOutput(r.Format, r.NoColor, r.PageSize)
	t.SetTitle("Members of " + r.TeamID)
	t.AppendHeader(table.Row{"Role", "UserID", "Email"})

	for _, member := range *team.Members {
		t.AppendRow([]interface{}{*member.Role, *member.UserId, *member.User.Email})
	}
	renderTable(t, r.Format)
	return nil
}

func (r *deleteMember) Execute([]string) error {
	if !r.YesIAmSure {
		if !verifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client, ctx, cancel := newUserAPIClient()
	defer cancel()

	team, res, err := client.TeamsApi.DeleteMember(ctx, r.TeamID, r.UserID).Execute()
	if err != nil {
		return apiError(res, err)
	}

	fmt.Printf("deleted member %s from %s\n", r.UserID, *team.TeamId)
	return nil
}
