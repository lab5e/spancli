package output

import (
	"encoding/json"
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type updateOutput struct {
	ID              commonopt.Collection
	NewCollectionID string `long:"new-collection-id" description:"new collection ID for output"`
	OID             oid
	Config          string `long:"config" description:"json configuration for output"`
	Tags            commonopt.Tags
	Enable          bool `long:"enable" description:"enable output"`
	Disable         bool `long:"disable" description:"disable output"`
	Type            otp
}

func (c *updateOutput) Execute([]string) error {
	//Ensure configuration is valid JSON first

	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	updateRequest := spanapi.UpdateOutputRequest{}
	if c.Config != "" {
		cfg := &spanapi.OutputConfig{}
		if err := json.Unmarshal([]byte(c.Config), cfg); err != nil {
			return fmt.Errorf("invalid json configuration: %v", err)
		}
		updateRequest.Config = cfg
	}
	if c.Type.Type != "" {
		ot := spanapi.OutputType(c.Type.Type)
		if !ot.IsValid() {
			return fmt.Errorf("invalid output type: %s", c.Type.Type)
		}
		updateRequest.Type = ot.Ptr()
	}
	if len(c.Tags.Tags) > 0 {
		updateRequest.Tags = c.Tags.AsMap()
	}
	if c.NewCollectionID != "" {
		updateRequest.CollectionId = spanapi.PtrString(c.NewCollectionID)
	}
	if c.Enable {
		updateRequest.Enabled = spanapi.PtrBool(true)
	}
	if c.Disable {
		updateRequest.Enabled = spanapi.PtrBool(false)
	}
	o, res, err := client.OutputsApi.UpdateOutput(ctx, c.ID.CollectionID, c.OID.OutputID).Body(updateRequest).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Printf("Updated output %s\n", o.GetOutputId())
	return nil
}
