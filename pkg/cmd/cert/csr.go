package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

// csrCert creates a local certificate and a certificate signing request for the certificate. The CSR can
// contain a lot of additional fields but this is mostly for demostration purposes so most fields use
// the defaults.
type csrCert struct {
	ID           commonopt.CollectionAndDevice
	Cert         string `long:"cert" description:"local certificate file" required:"yes" default:"cert.crt"`
	Key          string `long:"key" description:"local key file" required:"yes" default:"key.pem"`
	Email        string `long:"email" description:"email address for CSR" required:"yes"`
	Organization string `long:"org" description:"organization name for CSR" required:"yes"`
}

func (cc *csrCert) Execute([]string) error {
	fmt.Println("Generating private key...")
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Printf("Got error generating private key: %v", err)
		return err
	}

	fmt.Println("Public key generate, sending CSR request...")
	template := &x509.CertificateRequest{
		PublicKeyAlgorithm: x509.ECDSA,
		PublicKey:          &privateKey.PublicKey,
		Subject: pkix.Name{
			CommonName:   cc.Email,
			Organization: []string{cc.Organization},
		},
		EmailAddresses: []string{cc.Email},
	}

	// Create the CSR
	csr, err := x509.CreateCertificateRequest(rand.Reader, template, privateKey)
	if err != nil {
		fmt.Printf("Error generating CSR: %v\n", err)
		return err
	}

	csrBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr})

	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	csrResponse, res, err := client.CertificatesApi.SignCertificate(ctx, cc.ID.CollectionID).Body(
		spanapi.SignCertificateRequest{
			DeviceId: spanapi.PtrString(cc.ID.DeviceID),
			Csr:      spanapi.PtrString(base64.StdEncoding.EncodeToString(csrBytes)),
		}).Execute()
	if err != nil {
		return helpers.APIError(res, err)
	}

	fmt.Println("CSR request completed,")
	// The response has certificate and chain in separate base64 encoded fields. The fields are PEM encoded
	certBytes, err := base64.StdEncoding.DecodeString(csrResponse.GetCertificate())
	if err != nil {
		return fmt.Errorf("invalid base64 cert encoding: %v", err)
	}

	chainBytes, err := base64.StdEncoding.DecodeString(csrResponse.GetChain())
	if err != nil {
		return fmt.Errorf("invalid base64 chain encoding: %v", err)
	}

	certBytes = append(certBytes, chainBytes...)
	if err := writeFile(cc.Cert, certBytes); err != nil {
		return err
	}
	fmt.Printf("Wrote certificate to %s\n", cc.Cert)

	privateBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		fmt.Printf("Could not marshal private key: %v\n", err)
		return err
	}
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateBytes})
	if err := writeFile(cc.Key, keyBytes); err != nil {
		return err
	}

	fmt.Printf("Wrote private key to %s\n", cc.Key)
	return nil
}
