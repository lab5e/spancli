package gateway

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listGateways struct {
	ID     commonopt.Collection
	Format commonopt.ListFormat
}

func (l *listGateways) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	list, res, err := client.GatewaysApi.ListGateways(ctx, l.ID.CollectionID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	t := helpers.NewTableOutput(l.Format)
	t.SetTitle("Gateways for collection %s", l.ID.CollectionID)

	t.AppendHeader(table.Row{
		"Gateway ID",
		"Name",
		"Built In",
		"Type",
		"Status",
		"Tags",
	})

	for _, gw := range list.Gateways {
		t.AppendRow(table.Row{
			gw.GetGatewayId(),
			gw.GetName(),
			gw.GetBuiltIn(),
			gw.GetType(),
			gw.GetStatus(),
			helpers.TagsToString(gw.GetTags()),
		})
	}
	t.AppendFooter(table.Row{fmt.Sprintf("%d gateway(s)", len(list.Gateways))})
	helpers.RenderTable(t, l.Format.Format)
	return nil
}
