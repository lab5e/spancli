package invite

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteInvite struct {
	TeamID     string `long:"team-id" description:"id of team we wish to delete invite from" required:"yes"`
	Code       string `long:"code" description:"invite code we wish to delete"`
	YesIAmSure bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
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
