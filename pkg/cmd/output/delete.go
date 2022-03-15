package output

import (
	"errors"
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type deleteOutput struct {
	ID     commonopt.Collection
	OID    oid
	Prompt commonopt.NoPrompt
}

func (c *deleteOutput) Execute([]string) error {
	if !c.Prompt.Check() {
		return errors.New("aborted by user")
	}

	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	o, res, err := client.OutputsApi.DeleteOutput(ctx, c.ID.CollectionID, c.OID.OutputID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("Removed output %s\n", o.GetOutputId())
	return nil
}
