package collection

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteCollection struct {
	ID     commonopt.Collection
	Prompt commonopt.NoPrompt
}

func (r *deleteCollection) Execute([]string) error {
	if !r.Prompt.Check() {
		return fmt.Errorf("user aborted delete")
	}

	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	col, res, err := client.CollectionsApi.DeleteCollection(ctx, r.ID.CollectionID).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("deleted collection %s\n", *col.CollectionId)
	return nil
}
