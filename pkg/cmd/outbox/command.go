package outbox

type Command struct {
	Add    addOutboxCmd    `command:"add" description:"add message to oubox"`
	List   listOutboxCmd   `command:"list" description:"list messages in outbox"`
	Delete deleteOutboxCmd `command:"delete" description:"delete message from outbox"`
	Watch  watchOutboxCmd  `command:"watch" description:"monitor message in outbox"`
}
