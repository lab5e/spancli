package cert

import (
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type validateCert struct {
	ID   commonopt.CollectionAndDevice
	Cert string `long:"cert" description:"local certificate file" required:"yes" default:"cert.crt"`
}

func (vc *validateCert) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	fmt.Printf("Validating cert in %s\n", vc.Cert)
	certBytes, err := os.ReadFile(vc.Cert)
	if err != nil {
		return err
	}
	// Pull just the certificate if this file contains multiple certificates
	block, _ := pem.Decode(certBytes)

	ver, res, err := client.CertificatesApi.VerifyCertificate(ctx, vc.ID.CollectionID).Body(
		spanapi.VerifyCertificateRequest{
			DeviceId:    spanapi.PtrString(vc.ID.DeviceID),
			Certificate: spanapi.PtrString(base64.StdEncoding.EncodeToString(pem.EncodeToMemory(block))),
		}).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	if !ver.GetValid() {
		fmt.Println("**** Invalid certificate!")
		for _, errStr := range ver.Errors {
			fmt.Println("    ", errStr)
		}
		return errors.New("invalid certificate")
	}
	fmt.Println("Certificate is OK")
	return nil
}
