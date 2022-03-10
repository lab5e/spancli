package collection

import (
	"encoding/json"
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
)

type getCollection struct {
	CollectionID string `long:"collection-id" description:"Span collection ID" required:"yes"`
}

func (r *getCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	col, res, err := client.CollectionsApi.RetrieveCollection(ctx, r.CollectionID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	jsonData, err := json.MarshalIndent(col, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}
