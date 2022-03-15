package output

type Command struct {
	Add        addOutput     `command:"add" description:"add output"`
	ConfigHelp configHelp    `command:"config" description:"show configuration options"`
	List       listOutput    `command:"list" description:"list outputs"`
	Delete     deleteOutput  `command:"delete" description:"delete output"`
	Update     updateOutput  `command:"update" description:"update output"`
	Enable     enableOutput  `command:"enable" description:"enable output"`
	Disable    disableOutput `command:"disable" description:"disable output"`
}
