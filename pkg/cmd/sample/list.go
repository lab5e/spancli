package sample

import "github.com/lab5e/spancli/pkg/commonopt"

type listSamples struct {
	Format  commonopt.ListFormat
	Keyword string `long:"keyword" description:"Filter on keyword (c, golang, dart, swift, rust, zephyr, arduino, esp32...)" `
}

// Execute runs the version command
func (l *listSamples) Execute([]string) error {
	return ListSamples(l.Format, l.Keyword)
}
