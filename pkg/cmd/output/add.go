package output

import (
	"encoding/json"
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type addOutput struct {
	ID       commonopt.Collection
	Config   string `long:"config" description:"json configuration for output" required:"yes"`
	Tags     commonopt.Tags
	Disabled bool `long:"disabled" description:"disable output when creating it"`
	Type     otp
}

func (c *addOutput) Execute([]string) error {
	//Ensure configuration is valid JSON first
	cfg := &spanapi.OutputConfig{}
	if err := json.Unmarshal([]byte(c.Config), cfg); err != nil {
		return fmt.Errorf("invalid json configuration: %v", err)
	}

	if c.Type.Type == "" {
		return fmt.Errorf("must specify output type")
	}

	ot := spanapi.OutputType(c.Type.Type)
	if !ot.IsValid() {
		return fmt.Errorf("invalid output type: %s", c.Type.Type)
	}

	client, ctx, done := helpers.NewSpanAPIClient()
	defer done()

	o, res, err := client.OutputsApi.CreateOutput(ctx, c.ID.CollectionID).Body(
		spanapi.CreateOutputRequest{
			Type:    ot.Ptr(),
			Config:  cfg,
			Enabled: spanapi.PtrBool(!c.Disabled),
			Tags:    c.Tags.AsMap(),
		},
	).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}
	fmt.Printf("Created output %s on collection %s\n", o.GetOutputId(), c.ID.CollectionID)
	return nil
}
