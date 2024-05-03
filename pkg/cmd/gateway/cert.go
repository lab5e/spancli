package gateway

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type gatewayCerts struct {
	ID     commonopt.CollectionAndGateway
	Format commonopt.ListFormat
}

func (l *gatewayCerts) Execute([]string) error {
	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	list, res, err := client.GatewaysApi.GatewayCertificates(ctx, l.ID.CollectionID, l.ID.GatewayID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	t := helpers.NewTableOutput(l.Format)
	t.SetTitle("Certificates for gateway %s", l.ID.GatewayID)

	t.AppendHeader(table.Row{
		"Serial #",
		"Expires",
	})

	for _, gw := range list.Certificates {
		ns, err := strconv.ParseInt(gw.GetExpires(), 10, 64)
		if err != nil {
			return errors.New("invalid expire time for certificate")
		}
		ts := time.UnixMilli(ns)
		t.AppendRow(table.Row{
			gw.GetCertificateSerial(),
			ts.Format(time.RFC3339),
		})
	}
	t.AppendFooter(table.Row{fmt.Sprintf("%d certificates(s)", len(list.Certificates))})
	helpers.RenderTable(t, l.Format.Format)

	return nil
}
