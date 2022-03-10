package outbox

type Command struct {
	CollectionID string          `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	DeviceID     string          `long:"device-id" description:"device id"`
	Add          addOutboxCmd    `command:"add" description:"add message to oubox"`
	List         listOutboxCmd   `command:"list" description:"list messages in outbox"`
	Delete       deleteOutboxCmd `command:"delete" description:"delete message from outbox"`
	Watch        watchOutboxCmd  `command:"watch" description:"monitor message in outbox"`
}
