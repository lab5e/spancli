package team

import (
	"fmt"

	"github.com/lab5e/go-userapi"
	"github.com/lab5e/spancli/pkg/helpers"
)

type addTeam struct {
	Name string `long:"name" description:"team name"`
}

func (r *addTeam) Execute([]string) error {
	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	team := userapi.Team{
		Tags: &map[string]string{},
	}

	if r.Name != "" {
		(*team.Tags)["name"] = r.Name
	}

	team, res, err := client.TeamsApi.CreateTeam(ctx).Body(team).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("created team %s\n", *team.TeamId)
	return nil
}
