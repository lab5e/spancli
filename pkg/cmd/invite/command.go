package invite

// Command holds the subcommands for the invite command
type Command struct {
	Add    addInvite    `command:"add" description:"add invite for team"`
	List   listInvite   `command:"list" alias:"ls" description:"list invites for team"`
	Delete deleteInvite `command:"delete" alias:"del" description:"delete invite from team"`
	Accept acceptInvite `command:"accept" description:"accept invite"`
}
