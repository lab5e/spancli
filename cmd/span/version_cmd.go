package main

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/helpers"
	"github.com/lab5e/spanclient-go/v4"
)

type versionCmd struct {
	// No parameters
}

func (v *versionCmd) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	info, _, err := client.SystemApi.GetSystemInfo(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Server reports version %s (%s)\n", info.Version, info.ReleaseName)
	fmt.Printf("This utility assumes version %s\n", helpers.ExpectedVersion)
	return nil
}
