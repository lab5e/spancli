package sample

import (
	"fmt"
)

// Command is the sample command. Samples are hosted on GitHub and the topics (aka tags)
// for repositories are used to filter the samples. The current client is limited
// to a maximum of 100 repositories
type Command struct {
	// No parameters
	List listSamples `command:"list" alias:"ls" description:"List available samples"`
}

// Execute runs the version command
func (*Command) Execute([]string) error {
	fmt.Println("Sample command")
	return nil
}
