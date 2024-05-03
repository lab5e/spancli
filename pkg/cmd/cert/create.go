package cert

import (
	"encoding/base64"
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

type createCert struct {
	ID   commonopt.CollectionAndDeviceOrGateway
	Cert string `long:"cert" description:"certificate file name" default:"cert.crt" required:"yes"`
	Key  string `long:"key" description:"private key file" default:"key.pem" required:"yes"`
}

func (cc *createCert) Execute([]string) error {
	if !cc.ID.Valid() {
		return fmt.Errorf("either device or gateway must be specified")
	}
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	var deviceIDptr *string = nil
	var gatewayIDptr *string = nil
	if cc.ID.DeviceID != "" {
		deviceIDptr = &cc.ID.DeviceID
	}
	if cc.ID.GatewayID != "" {
		gatewayIDptr = &cc.ID.GatewayID
	}
	cert, res, err := client.CertificatesApi.CreateCertificate(ctx, cc.ID.CollectionID).Body(
		spanapi.CreateCertificateRequest{
			DeviceId:  deviceIDptr,
			GatewayId: gatewayIDptr,
		}).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	// The response has certificate and chain in separate base64 encoded fields. The fields are PEM encoded
	certBytes, err := base64.StdEncoding.DecodeString(cert.GetCertificate())
	if err != nil {
		return fmt.Errorf("invalid base64 cert encoding: %v", err)
	}

	chainBytes, err := base64.StdEncoding.DecodeString(cert.GetChain())
	if err != nil {
		return fmt.Errorf("invalid base64 chain encoding: %v", err)
	}

	certBytes = append(certBytes, chainBytes...)
	if err := writeFile(cc.Cert, certBytes); err != nil {
		return err
	}
	fmt.Printf("Wrote certificate to %s\n", cc.Cert)

	keyBytes, err := base64.StdEncoding.DecodeString(cert.GetPrivateKey())
	if err != nil {
		return fmt.Errorf("invalid base64 key encoding: %v", err)
	}

	if err := writeFile(cc.Key, keyBytes); err != nil {
		return err
	}
	fmt.Printf("Wrote key to %s\n", cc.Key)

	return nil
}
