package outbox

import "fmt"

type listOutboxCmd struct {
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

func (*listOutboxCmd) Execute([]string) error {
	fmt.Println("list outbox")
	return nil
}
