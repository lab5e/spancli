package gateway

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

// addGateway is a command to add a new gateway. The type of gateway is
// always set to "user" (the built-ins cannot be created by clients for
// obvious reasons)
type addGateway struct {
	ID     commonopt.Collection
	Name   string   `long:"name" description:"name of gateway" required:"yes"`
	Config []string `long:"config" description:"configuration parameter"`
	Tags   commonopt.Tags
}

func (a *addGateway) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	gwType := spanapi.GATEWAYTYPE_CUSTOM
	create := spanapi.CreateGatewayRequest{
		Type: &gwType,
		Name: spanapi.PtrString(a.Name),
		Tags: a.Tags.AsMap(),
		Config: &spanapi.GatewayConfig{
			User: &spanapi.GatewayCustomConfig{
				Params: helpers.AsMap(a.Config),
			},
		},
	}

	gw, res, err := client.GatewaysApi.CreateGateway(ctx, a.ID.CollectionID).Body(create).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("Gateway %s created\n", gw.GetGatewayId())
	return nil
}
