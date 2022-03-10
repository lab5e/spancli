package invite

import (
	"fmt"

	"github.com/lab5e/go-userapi"
	"github.com/lab5e/spancli/pkg/helpers"
)

type addInvite struct {
	TeamID string `long:"team-id" description:"id of team we wish to add invite to" required:"yes"`
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
