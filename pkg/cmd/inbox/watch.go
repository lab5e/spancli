package inbox

import "fmt"

type watchInboxCmd struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id" env:"SPAN_DEVICE_ID"`
	//lint:ignore SA5008 Multiple choices makes the linter unhappy
	Format string `long:"format" default:"text" description:"which output format to use" choice:"text" choice:"json"`
}

func (*watchInboxCmd) Execute([]string) error {
	fmt.Println("Watch inbox")
	return nil
}
