package collection

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type addCollection struct {
	TeamID string `long:"team-id" description:"team the collection belongs to"`
	Tags   commonopt.Tags
}

func (r *addCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	collection := spanapi.CreateCollectionRequest{
		Tags: r.Tags.AsMap(),
	}
	// Team ID is optional; the private team will be used if omitted

	if r.TeamID != "" {
		collection.TeamId = &r.TeamID
	}

	col, res, err := client.CollectionsApi.CreateCollection(ctx).Body(collection).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}
	fmt.Printf("created collection %s\n", col.GetCollectionId())
	return nil
}
