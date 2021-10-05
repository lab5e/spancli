package main

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/go-userapi"
)

type teamCmd struct {
	Add    addTeam    `command:"add" description:"create new team"`
	List   listTeams  `command:"list" alias:"ls" description:"list teams"`
	Delete deleteTeam `command:"delete" alias:"del" description:"delete team"`
}

type addTeam struct {
	Name string `long:"name" description:"team name"`
}

type listTeams struct {
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

type deleteTeam struct {
	TeamID     string `long:"team-id" description:"id of team we wish to delete" required:"yes"`
	YesIAmSure bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

func (r *addTeam) Execute([]string) error {
	client, ctx, cancel := newUserAPIClient()
	defer cancel()

	team := userapi.Team{
		Tags: &map[string]string{},
	}

	if r.Name != "" {
		(*team.Tags)["name"] = r.Name
	}

	team, res, err := client.TeamsApi.CreateTeam(ctx).Body(team).Execute()
	if err != nil {
		return apiError(res, err)
	}

	fmt.Printf("created team '%s'\n", *team.TeamId)
	return nil
}

func (r *listTeams) Execute([]string) error {
	client, ctx, cancel := newUserAPIClient()
	defer cancel()

	teamList, res, err := client.TeamsApi.ListTeams(ctx).Execute()
	if err != nil {
		return apiError(res, err)
	}

	if r.Format == "json" {
		json, err := json.MarshalIndent(*teamList.Teams, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := newTableOutput(r.Format, r.NoColor, r.PageSize)
	t.SetTitle("Teams")
	t.AppendHeader(table.Row{"ID", "Name"})

	for _, team := range *teamList.Teams {
		t.AppendRow([]interface{}{*team.TeamId, team.GetTags()["name"]})
	}
	renderTable(t, r.Format)

	return nil
}

func (r *deleteTeam) Execute([]string) error {
	client, ctx, cancel := newUserAPIClient()
	defer cancel()

	if !r.YesIAmSure {
		if !verifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	team, res, err := client.TeamsApi.DeleteTeam(ctx, r.TeamID).Execute()
	if err != nil {
		return apiError(res, err)
	}

	fmt.Printf("deleted team '%s'\n", *team.TeamId)

	return nil
}
