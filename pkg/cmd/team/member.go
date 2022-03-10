package team

type memberCmd struct {
	List   listMembers  `command:"list" alias:"ls" description:""`
	Delete deleteMember `command:"delete" alias:"del" description:""`
}
