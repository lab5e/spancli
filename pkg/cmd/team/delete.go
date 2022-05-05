package team

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteTeam struct {
	TeamID string `long:"team-id" description:"id of team we wish to delete" required:"yes"`
	Prompt commonopt.NoPrompt
}

func (r *deleteTeam) Execute([]string) error {
	if !r.Prompt.Check() {
		return fmt.Errorf("user aborted delete")
	}

	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	team, res, err := client.TeamsApi.DeleteTeam(ctx, r.TeamID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("deleted team %s\n", *team.TeamId)

	return nil
}
