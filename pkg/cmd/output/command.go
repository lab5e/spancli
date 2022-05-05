package output

// Command holds the subcommands for the output command
type Command struct {
	Add        addOutput    `command:"add" description:"add output"`
	ConfigHelp configHelp   `command:"config" alias:"cfg" description:"show configuration options"`
	List       listOutput   `command:"list" alias:"ls" description:"list outputs"`
	Update     updateOutput `command:"update" alias:"up" description:"update output"`
	Logs       outputLogs   `command:"logs" alias:"lg" description:"show output logs"`
	Delete     deleteOutput `command:"delete" alias:"del" description:"delete output"`
	Status     outputStatus `command:"status" alias:"s" description:"show status"`
	Get        getOutput    `command:"get" description:"get output details"`
}
