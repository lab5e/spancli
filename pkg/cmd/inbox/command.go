package inbox

import "fmt"

type Command struct {
	CollectionID string        `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string        `long:"device-id" description:"device id"`
	List         listInboxCmd  `command:"list" description:"list contents of inbox"`
	Watch        watchInboxCmd `command:"watch" description:"watch contents of inbox"`
}

type listInboxCmd struct {
	Limit    int32  `long:"limit" description:"max number of entries to fetch" default:"30"`
	Start    string `long:"start" description:"start of time range in milliseconds since epoch"`
	End      string `long:"end" description:"end of time range in milliseconds since epoch"`
	Decode   bool   `long:"decode" description:"decode payload"`
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

func (*listInboxCmd) Execute([]string) error {
	fmt.Println("List inbox")
	return nil
}

type watchInboxCmd struct {
	Format string `long:"format" default:"text" description:"which output format to use" choice:"text" choice:"json"`
}

func (*watchInboxCmd) Execute([]string) error {
	fmt.Println("Watch inbox")
	return nil
}
