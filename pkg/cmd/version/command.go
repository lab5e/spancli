package version

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
)

type Command struct {
	// No parameters
}

func (*Command) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	info, _, err := client.SystemApi.GetSystemInfo(ctx).Execute()
	if err != nil {
		return err
	}

	fmt.Printf("Server reports version %s (%s)\n", *info.Version, *info.ReleaseName)
	fmt.Printf("This utility assumes version %s\n", helpers.ExpectedVersion)
	return nil
}
