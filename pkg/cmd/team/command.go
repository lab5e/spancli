package team

type Command struct {
	Add    addTeam    `command:"add" description:"create new team"`
	Get    getTeam    `command:"get" description:"get team details"`
	List   listTeams  `command:"list" alias:"ls" description:"list teams"`
	Delete deleteTeam `command:"delete" alias:"del" description:"delete team"`

	Members memberCmd `command:"member" description:"manage team members"`
}
