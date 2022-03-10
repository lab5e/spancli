package inbox

type Command struct {
	List  listInboxCmd  `command:"list" description:"list contents of inbox"`
	Watch watchInboxCmd `command:"watch" description:"watch contents of inbox"`
}
