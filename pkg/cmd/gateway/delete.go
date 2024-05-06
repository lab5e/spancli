package gateway

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteGateway struct {
	ID     commonopt.CollectionAndGateway
	Prompt commonopt.NoPrompt
}

func (d *deleteGateway) Execute([]string) error {
	if !d.Prompt.Check() {
		return fmt.Errorf("delete aborted")
	}

	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	gw, res, err := client.GatewaysApi.DeleteGateway(ctx, d.ID.CollectionID, d.ID.GatewayID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("Deleted gateway %s\n", gw.GetGatewayId())
	return nil
}
