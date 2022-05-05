package team

import (
	"encoding/json"
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
)

type getTeam struct {
	TeamID string `long:"team-id" description:"id of team" required:"yes"`
}

func (r *getTeam) Execute([]string) error {
	client, ctx, cancel := helpers.NewUserAPIClient()
	defer cancel()

	team, res, err := client.TeamsApi.RetrieveTeam(ctx, r.TeamID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	jsonData, err := json.MarshalIndent(team, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}
