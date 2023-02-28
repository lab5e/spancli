package blob

// Command holds the subcommands for the cert command
type Command struct {
	Get  downloadBlob `command:"get" description:"Download blob"`
	List listBlobs    `command:"list" alias:"ls" description:"list blobs"`
}
