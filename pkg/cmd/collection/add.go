package collection

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

type addCollection struct {
	TeamID string   `long:"team-id" description:"team the collection belongs to" required:"yes"`
	Tags   []string `long:"tag" description:"Set tag value (name:value)"`
	Name   string   `long:"name" description:"name of the collection"`
}

func (r *addCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	collection := spanapi.CreateCollectionRequest{
		TeamId: &r.TeamID,
		Tags:   helpers.TagMerge(nil, r.Tags),
	}

	if r.Name != "" {
		(*collection.Tags)["name"] = r.Name
	}

	col, res, err := client.CollectionsApi.CreateCollection(ctx).Body(collection).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("created collection %s\n", *col.CollectionId)
	return nil
}
