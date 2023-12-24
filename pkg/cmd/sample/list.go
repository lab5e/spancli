package sample

import "github.com/lab5e/spancli/pkg/commonopt"

type listSamples struct {
	Format commonopt.ListFormat
}

// Execute runs the version command
func (l *listSamples) Execute([]string) error {
	return ListSamples(l.Format)
}
