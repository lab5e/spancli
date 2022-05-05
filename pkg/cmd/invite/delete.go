package invite

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteInvite struct {
	TeamID string `long:"team-id" description:"id of team we wish to delete invite from" required:"yes"`
	Code   string `long:"code" description:"invite code we wish to delete"`
	Prompt commonopt.NoPrompt
}

func (r *deleteInvite) Execute([]string) error {
	if !r.Prompt.Check() {
		return fmt.Errorf("user aborted delete")
	}

	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	resp, res, err := client.TeamsApi.DeleteInvite(ctx, r.TeamID, r.Code).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("deleted invite %s from %s\n", *resp.Invite.Code, r.TeamID)

	return nil
}
