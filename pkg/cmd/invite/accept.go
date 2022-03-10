package invite

import (
	"fmt"

	"github.com/lab5e/go-userapi"
	"github.com/lab5e/spancli/pkg/helpers"
)

type acceptInvite struct {
	Code string `long:"code" description:"invite code" required:"yes"`
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
