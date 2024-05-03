package gateway

// Command holds the subcommands for the gateway command
type Command struct {
	List   listGateways  `command:"list" alias:"ls" description:"list gateways"`
	Add    addGateway    `command:"add" description:"create gateway"`
	Delete deleteGateway `command:"delete" alias:"rm" description:"remove gateway"`
	Update updateGateway `command:"update" description:"update gateway"`
	Watch  watchGateway  `command:"watch" alias:"monitor" description:"monitor gateway activity"`
	Cert   gatewayCerts  `command:"cert" alias:"certificates" description:"list gateway certificates"`
}
