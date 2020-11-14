package main

type deviceCmd struct {
	Add    addDevice    `command:"add" description:"create device"`
	Get    getDevice    `command:"get" description:"get device"`
	List   listDevice   `command:"list" description:"list devices"`
	Delete deleteDevice `command:"delete" description:"delete device"`
}

type addDevice struct {
	CollectionID string
	Name         string
	IMSI         string
	IMEI         string
}

type getDevice struct{}

type listDevice struct{}

type deleteDevice struct{}
