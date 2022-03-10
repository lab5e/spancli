package collection

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

type updateCollection struct {
	CollectionID string   `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	Tags         []string `long:"tag" description:"Set tag value (name:value)"`
}

func (u *updateCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	update := spanapi.UpdateCollectionRequest{
		Tags: helpers.TagMerge(nil, u.Tags),
	}

	collectionUpdated, res, err := client.CollectionsApi.UpdateCollection(ctx, u.CollectionID).Body(update).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("updated collection %s\n", *collectionUpdated.CollectionId)
	return nil
}
