package firmware

// Command holds the subcommands for the firmware command
type Command struct {
	List    listFirmware    `command:"list" alias:"ls" description:"list firmware images"`
	Upload  uploadFirmware  `command:"upload" alias:"up" description:"upload firmware image"`
	Update  updateFirmware  `command:"update" description:"update firmware image"`
	Delete  deleteFirmware  `command:"delete" alias:"del" description:"delete firmware image"`
	Monitor monitorFirmware `command:"monitor" alias:"mon" description:"monitor firmware update status"`
	Reset   resetError      `command:"clear" description:"clear firmware errors on device"`
	Get     getFirmware     `command:"get" description:"get firmware details"`
}
