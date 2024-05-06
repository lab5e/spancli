package gateway

import (
	"fmt"
	"regexp"

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

	gwType := spanapi.GatewayType("user")
	create := spanapi.CreateGatewayRequest{
		Type:   &gwType,
		Name:   spanapi.PtrString(a.Name),
		Tags:   a.Tags.AsMap(),
		Config: a.asConfigMap(),
	}

	gw, res, err := client.GatewaysApi.CreateGateway(ctx, a.ID.CollectionID).Body(create).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("Gateway %s created successfully", gw.GetGatewayId())
	return nil
}

// configRegex matches tags of the form:
//
//	foo:bar
//	foo:"bar baz"
//	foo:
//	foo:""
var configRegex = regexp.MustCompile(`^\s*(\S+):("?)(.*?)("?)\s*$`)

func (a *addGateway) asConfigMap() *spanapi.GatewayConfig {
	params := make(map[string]string)

	for _, elt := range a.Config {
		res := configRegex.FindStringSubmatch(elt)
		if len(res) != 5 {
			continue
		}
		params[res[1]] = res[3]
	}
	return &spanapi.GatewayConfig{
		User: &spanapi.GatewayCustomConfig{
			Params: &params,
		},
	}
}
