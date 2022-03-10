package collection

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteCollection struct {
	CollectionID string `long:"collection-id" description:"Span collection ID" required:"yes"`
	YesIAmSure   bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

func (r *deleteCollection) Execute([]string) error {
	if !r.YesIAmSure {
		if !helpers.VerifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	col, res, err := client.CollectionsApi.DeleteCollection(ctx, r.CollectionID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("deleted collection %s\n", *col.CollectionId)
	return nil
}
