package main

type outputCmd struct {
	Add     addOutput     `command:"add" description:"add output"`
	List    listOutput    `command:"list" description:"list outputs"`
	Delete  deleteOutput  `command:"delete" description:"delete output"`
	Update  updateOutput  `command:"update" description:"update output"`
	Enable  enableOutput  `command:"enable" description:"enable output"`
	Disable disableOutput `command:"disable" description:"disable output"`
}

type addOutput struct {
}

type listOutput struct {
}

type deleteOutput struct {
}

type updateOutput struct {
}

type enableOutput struct {
}

type disableOutput struct {
}
