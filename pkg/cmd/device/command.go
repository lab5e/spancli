package device

type Command struct {
	Add    addDevice    `command:"add" description:"create device"`
	Get    getDevice    `command:"get" description:"get device"`
	Update updateDevice `command:"update" description:"update device"`
	List   listDevices  `command:"list" alias:"ls" description:"list devices"`
	Delete deleteDevice `command:"delete" alias:"del" description:"delete device"`
}
