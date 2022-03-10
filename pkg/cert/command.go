package cert

type Command struct {
	Create createCert `command:"create" description:"create a new client certificate"`
	Sign   signCert   `command:"sign" description:"sign a certificate"`
	Check  checkCert  `command:"check" description:"verify and display a certificate"`
}
