package blob

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/global"
	"github.com/lab5e/spancli/pkg/helpers"
)

type downloadBlob struct {
	OutputFile string `long:"output" description:"Output file name for blob" required:"yes"`
	ID         commonopt.Collection
	BlobID     string `long:"blob-id" description:"Blob ID" required:"yes"`
}

func (r *downloadBlob) Execute([]string) error {
	// Retrieve the blob, fix the path and download it
	defaultHost := "https://api.lab5e.com/"
	if global.Options.OverrideEndpoint != "" {
		defaultHost = global.Options.OverrideEndpoint
	}

	u, err := url.Parse(defaultHost)
	if err != nil {
		fmt.Printf("Could not parse host name: %v", err)
	}
	u.Path = fmt.Sprintf("/span/collections/%s/blobs/%s", r.ID.CollectionID, r.BlobID)

	fmt.Printf("Downloading blob %s...\n", r.BlobID)

	// Open output file for writing
	_, err = os.Stat(r.OutputFile)
	if !os.IsNotExist(err) {
		return errors.New("output file already exists")
	}
	outFile, err := os.Create(r.OutputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	buf := bytes.Buffer{}
	req, err := http.NewRequest("GET", u.String(), &buf)
	if err != nil {
		fmt.Println("Could not create request")
		return err
	}
	req.Header.Add("X-API-Token", global.Options.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return helpers.APIError(res, errors.New("expected 200 OK response"))
	}

	written, err := io.Copy(outFile, res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Wrote %d bytes to %s\n", written, r.OutputFile)

	return nil
}
