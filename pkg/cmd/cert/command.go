package cert

type Command struct {
	Create    createCert   `command:"create" description:"create a new client certificate"`
	Convert   convertCert  `command:"convert" description:"convert to other formats"`
	CreateCSR csrCert      `command:"csr" description:"create a local certificate and a certificate signing request"`
	Validate  validateCert `command:"validate" description:"verify and display a certificate"`
}
