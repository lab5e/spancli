package inbox

import "fmt"

type listInboxCmd struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id" env:"SPAN_DEVICE_ID"`

	Limit  int32  `long:"limit" description:"max number of entries to fetch" default:"30"`
	Start  string `long:"start" description:"start of time range in milliseconds since epoch"`
	End    string `long:"end" description:"end of time range in milliseconds since epoch"`
	Decode bool   `long:"decode" description:"decode payload"`

	//lint:ignore SA5008 Multiple choices makes the linter unhappy
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

func (*listInboxCmd) Execute([]string) error {
	fmt.Println("List inbox")
	return nil
}
