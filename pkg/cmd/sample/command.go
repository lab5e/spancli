package sample

// Command is the sample command. Samples are hosted on GitHub and the topics (aka tags)
// for repositories are used to filter the samples. The current client is limited
// to a maximum of 100 repositories
type Command struct {
	// No parameters
	List   listSamples  `command:"list" alias:"ls" description:"List available samples"`
	Create createSample `command:"create" alias:"new" description:"Create a sample"`
}
