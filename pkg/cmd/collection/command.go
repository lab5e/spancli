package collection

type Command struct {
	Add    addCollection    `command:"add" description:"create new collection"`
	Get    getCollection    `command:"get" description:"get collection details"`
	List   listCollection   `command:"list" alias:"ls" description:"list collections"`
	Delete deleteCollection `command:"delete" alias:"del" description:"delete collection"`
	Update updateCollection `command:"update" alias:"up" description:"update collection"`
}
