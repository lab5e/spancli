package invite

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/go-userapi"
	"github.com/lab5e/spancli/pkg/helpers"
)

type Command struct {
	Add    addInvite    `command:"add" description:"add invite for team"`
	List   listInvite   `command:"list" alias:"ls" description:"list invites for team"`
	Delete deleteInvite `command:"delete" alias:"del" description:"delete invite from team"`
	Accept acceptInvite `command:"accept" description:"accept invite"`
}

type addInvite struct {
	TeamID string `long:"team-id" description:"id of team we wish to add invite to" required:"yes"`
}

type listInvite struct {
	TeamID   string `long:"team-id" description:"id of team we wish to list invites for" required:"yes"`
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

type deleteInvite struct {
	TeamID     string `long:"team-id" description:"id of team we wish to delete invite from" required:"yes"`
	Code       string `long:"code" description:"invite code we wish to delete"`
	YesIAmSure bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

type acceptInvite struct {
	Code string `long:"code" description:"invite code" required:"yes"`
}

func (r *addInvite) Execute([]string) error {
	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	invite, res, err := client.TeamsApi.GenerateInvite(ctx, r.TeamID).Body(*userapi.NewInviteRequest()).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("created invite code for team %s: %s\n", r.TeamID, *invite.Code)
	return nil
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

func (r *deleteInvite) Execute([]string) error {
	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	if !r.YesIAmSure {
		if !helpers.VerifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	resp, res, err := client.TeamsApi.DeleteInvite(ctx, r.TeamID, r.Code).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("deleted invite %s from %s\n", *resp.Invite.Code, r.TeamID)

	return nil
}

func (r *acceptInvite) Execute([]string) error {
	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	resp, res, err := client.TeamsApi.AcceptInvite(ctx).Body(userapi.AcceptInviteRequest{Code: &r.Code}).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("accepted invite to team %s", *resp.TeamId)

	return nil
}
