package output

type Command struct {
	Add        addOutput    `command:"add" description:"add output"`
	ConfigHelp configHelp   `command:"config" description:"show configuration options"`
	List       listOutput   `command:"list" description:"list outputs"`
	Update     updateOutput `command:"update" description:"update output"`
	Logs       outputLogs   `command:"logs" description:"show output logs"`
	Delete     deleteOutput `command:"delete" description:"delete output"`
	Status     outputStatus `command:"status" description:"show status"`
	Get        getOutput    `command:"get" description:"get output details"`
}
