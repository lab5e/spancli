package gateway

// Command holds the subcommands for the gateway command
type Command struct {
	List         listGateways  `command:"list" alias:"ls" description:"list gateways"`
	Add          addGateway    `command:"add" description:"create gateway"`
	Delete       deleteGateway `command:"delete" alias:"rm" description:"remove gateway"`
	Update       updateGateway `command:"update" description:"update gateway"`
	Cert         gatewayCerts  `command:"cert" alias:"certificates" description:"list gateway certificates"`
	SampleConfig sampleConfigs `command:"sample-configs" description:"show sample configurations"`
}
