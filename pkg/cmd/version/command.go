package version

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
)

// Command is the version command
type Command struct {
	// No parameters
}

// Execute runs the version command
func (*Command) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	info, _, err := client.SpanApi.GetSystemInfo(ctx).Execute()
	if err != nil {
		return err
	}

	fmt.Printf("Server reports version %s (%s)\n", *info.Version, *info.ReleaseName)
	fmt.Printf("This utility assumes version %s\n", helpers.ExpectedVersion)
	return nil
}
