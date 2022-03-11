package cert

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

type printCert struct {
	Cert string `long:"cert" description:"local certificate file" required:"yes" default:"cert.crt"`
}

func (pc *printCert) Execute([]string) error {
	certBytes, err := os.ReadFile(pc.Cert)
	if err != nil {
		return err
	}
	fmt.Printf("%d bytes read from %s\n", len(certBytes), pc.Cert)

	block, remain := pem.Decode(certBytes)
	for block != nil {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}
		print(cert)
		block, remain = pem.Decode(remain)
	}
	// Read pem blocks and
	return nil
}

func print(cert *x509.Certificate) {
	fmt.Println("------------------------------------")
	fmt.Printf("Algorithm:           %s (%d bits)\n", cert.PublicKeyAlgorithm.String(), keySize(cert.PublicKey))
	fmt.Printf("Signature Algorithm: %s\n", cert.SignatureAlgorithm.String())
	fmt.Printf("Serial:              %s\n", hex.EncodeToString(cert.SerialNumber.Bytes()))
	fmt.Printf("Subject:             %s (issued by %s)\n", cert.Subject.CommonName, cert.Issuer.CommonName)
	if len(cert.Subject.Organization) > 0 {
		fmt.Printf("    Organization:    %s\n", cert.Subject.Organization[0])
	}
	if len(cert.Subject.OrganizationalUnit) > 0 {
		fmt.Printf("    Org. Unit:       %s\n", cert.Subject.OrganizationalUnit[0])
	}
	if len(cert.Subject.Country) > 0 {
		fmt.Printf("    Country:         %s\n", cert.Subject.Country[0])
	}
	if len(cert.EmailAddresses) > 0 {
		fmt.Printf("Email:               %s\n", cert.EmailAddresses[0])
	}
	fmt.Printf("CA:                  %t\n", cert.IsCA)
	fmt.Printf("Valid from:          %s\n", cert.NotBefore.Format(time.RFC3339))
	fmt.Printf("Valid to:            %s (%s)\n", cert.NotAfter.Format(time.RFC3339), cert.NotAfter.Sub(cert.NotBefore))
	fmt.Printf("Max path:            %d\n", cert.MaxPathLen)
}

// keySize returns the number of bits used for the key algorithm
func keySize(pub interface{}) int {
	switch k := pub.(type) {
	case *rsa.PublicKey:
		return k.Size() * 8
	case *ecdsa.PublicKey:
		return k.Curve.Params().BitSize
	case ed25519.PublicKey:
		return 32 * 8
	default:
		return 0
	}
}
