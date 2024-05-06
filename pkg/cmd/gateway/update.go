package gateway

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type updateGateway struct {
	ID     commonopt.CollectionAndGateway
	Tags   commonopt.Tags
	Name   string   `long:"name" description:"gateway name" required:"yes"`
	Config []string `long:"config" description:"configuration parameter"`
}

func (u *updateGateway) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	updateRequest := spanapi.UpdateGatewayRequest{}
	if u.Name != "" {
		updateRequest.Name = spanapi.PtrString(u.Name)
	}
	if len(u.Tags.Tags) > 0 {
		updateRequest.Tags = u.Tags.AsMap()
	}
	if len(u.Config) > 0 {
		updateRequest.Config = &spanapi.GatewayConfig{
			User: &spanapi.GatewayCustomConfig{
				Params: helpers.AsMap(u.Config),
			},
		}
	}
	gw, res, err := client.GatewaysApi.UpdateGateway(ctx, u.ID.CollectionID, u.ID.GatewayID).Body(updateRequest).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("Gateway %s is updated\n", gw.GetGatewayId())
	return nil
}
