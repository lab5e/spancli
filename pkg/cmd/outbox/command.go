package outbox

// Command holds the subcommands for the outbox command
type Command struct {
	Add    addOutboxCmd    `command:"add" description:"add message to oubox"`
	List   listOutboxCmd   `command:"list" alias:"ls" description:"list messages in outbox"`
	Delete deleteOutboxCmd `command:"delete" alias:"del" description:"delete message from outbox"`
	Watch  watchOutboxCmd  `command:"watch" alias:"w" description:"monitor message in outbox"`
}
