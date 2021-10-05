package main

import (
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

type versionCmd struct {
	// No parameters
}

func (v *versionCmd) Execute([]string) error {
	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	info, _, err := client.SystemApi.GetSystemInfo(ctx).Execute()
	if err != nil {
		return err
	}

	fmt.Printf("Server reports version %s (%s)\n", *info.Version, *info.ReleaseName)
	fmt.Printf("This utility assumes version %s\n", helpers.ExpectedVersion)
	return nil
}
