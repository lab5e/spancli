package cert

import (
	"fmt"
	"os"
)

type convertCert struct {
	In  string `long:"in" required:"yes" description:"input file, PEM encoded"`
	Out string `long:"out" required:"yes" description:"output file"`
	//lint:ignore SA5008 Multiple choice tags makes the linter unhappy
	Format string `long:"format" default:"der" description:"output format" choice:"der" choice:"hex-list" choice:"decimal-list" choice:"base64"`
}

func (cc *convertCert) Execute([]string) error {
	buf, err := os.ReadFile(cc.In)
	if err != nil {
		return err
	}

	outBytes, err := formatBytes(buf, cc.Format)
	if err != nil {
		return err
	}

	if err := writeFile(cc.Out, outBytes); err != nil {
		return err
	}

	fmt.Printf("Wrote %d bytes to %s\n", len(outBytes), cc.Out)
	return nil
}
