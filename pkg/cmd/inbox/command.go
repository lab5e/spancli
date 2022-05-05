package inbox

// Command holds the subcommands for the inbox command
type Command struct {
	List  listInboxCmd  `command:"list" alias:"ls" description:"list contents of inbox"`
	Watch watchInboxCmd `command:"watch" alias:"w" description:"watch contents of inbox"`
}
