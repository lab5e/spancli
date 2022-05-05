package commonopt

// QueryOptions holds common options for queries
type QueryOptions struct {
	Limit  int32  `long:"limit" description:"max number of entries to fetch" default:"30"`
	Start  string `long:"start" description:"start of time range in milliseconds since epoch"`
	End    string `long:"end" description:"end of time range in milliseconds since epoch"`
	Decode bool   `long:"decode" description:"decode payload"`
}
