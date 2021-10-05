package helpers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/lab5e/go-spanapi/v4"
)

// ExpectedVersion is the expected API version we'll use. This is for informational
// purposes only
const ExpectedVersion = "4.1.16"

// CheckVersion checks the server version of the Span API and emits a warning
// if there's a version mismatch.
func CheckVersion(ctx context.Context, client *spanapi.APIClient) {
	info, _, err := client.SystemApi.GetSystemInfo(ctx).Execute()
	if err != nil {
		fmt.Printf("Error checking version: %v\n", err)
		return
	}
	emitWarning := false

	if *info.Version < ExpectedVersion {
		// Emit warnings always when the expected version is GREATER than the
		// API version (we're pretty much certain that the features aren't
		// in line here)
		emitWarning = true
	}
	versionServer := strings.Split(*info.Version, ".")
	versionCli := strings.Split(ExpectedVersion, ".")

	if versionServer[0] != versionCli[0] {
		// Emit warning if the major version number doesn't match. There will be
		// breaking changes.
		emitWarning = true
	}
	if emitWarning {
		fmt.Fprintf(os.Stderr, `
		*** Warning:
			Server reports version %s but this client was built with version %s.
			The client might not work as intended. It's recommended to upgrade
			to a more current version.
		`, *info.Version, ExpectedVersion)
	}
}
