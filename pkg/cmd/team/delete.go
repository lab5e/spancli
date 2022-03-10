package team

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteTeam struct {
	TeamID     string `long:"team-id" description:"id of team we wish to delete" required:"yes"`
	YesIAmSure bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

func (r *deleteTeam) Execute([]string) error {
	if !r.YesIAmSure {
		if !helpers.VerifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	team, res, err := client.TeamsApi.DeleteTeam(ctx, r.TeamID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("deleted team %s\n", *team.TeamId)

	return nil
}
