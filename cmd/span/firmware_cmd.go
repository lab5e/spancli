package main

type firmwareCmd struct {
	Upload  uploadFirmware  `command:"upload" alias:"up" description:"upload firmware image"`
	Update  updateFirmware  `command:"update" description:"update firmware image"`
	Delete  deleteFirmware  `command:"delete" description:"delete firmware image"`
	Monitor monitorFirmware `command:"monitor" alias:"mon" description:"monitor firmware update status"`
}

type uploadFirmware struct {
}

type updateFirmware struct {
}

type deleteFirmware struct {
}

type monitorFirmware struct {
}
