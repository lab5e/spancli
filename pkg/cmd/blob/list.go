package blob

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type listBlobs struct {
	ID     commonopt.Collection
	Format commonopt.ListFormat
}

func (r *listBlobs) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	resp, res, err := client.BlobsApi.ListBlobs(ctx, r.ID.CollectionID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	if resp.Blobs == nil {
		fmt.Printf("no blobs\n")
		return nil
	}

	if r.Format.Format == "json" {
		js, err := json.MarshalIndent(resp.Blobs, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(js))
		return nil
	}

	t := helpers.NewTableOutput(r.Format)
	t.SetTitle("Blobs in %s", r.ID.CollectionID)
	t.AppendHeader(table.Row{
		"ID",
		"Content-Type",
		"Size",
		"Uploaded",
		"Path",
	})

	for _, blob := range resp.Blobs {
		ts, err := strconv.ParseInt(blob.GetCreated(), 10, 64)
		if err != nil {
			fmt.Printf("Error parsing created field\n")
			return err
		}

		t.AppendRow(table.Row{
			blob.GetBlobId(),
			blob.GetContentType(),
			blob.GetSize(),
			time.Unix(0, ts).Format(time.RFC3339),
			blob.GetBlobPath(),
		})
	}
	helpers.RenderTable(t, r.Format.Format)
	return nil
}
