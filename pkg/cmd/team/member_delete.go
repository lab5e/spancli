package team

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteMember struct {
	TeamID     string `long:"team-id" description:"id of team" required:"yes"`
	UserID     string `long:"user-id" description:"id of user we want to remove" required:"yes"`
	YesIAmSure bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

func (r *deleteMember) Execute([]string) error {
	if !r.YesIAmSure {
		if !helpers.VerifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	team, res, err := client.TeamsApi.DeleteMember(ctx, r.TeamID, r.UserID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("deleted member %s from %s\n", r.UserID, *team.TeamId)
	return nil
}