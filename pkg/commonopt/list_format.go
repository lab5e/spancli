package commonopt

type ListFormat struct {
	//lint:ignore SA5008 Multiple choices makes the linter unhappy
	Format          string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor         bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize        int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
	MaxPayloadWdith int    `long:"payload-width" description:"maximum width of payload" default:"50"`
	NumericDate     bool   `long:"numeric-date" description:"display dates as ms-since-epoch"`
}
