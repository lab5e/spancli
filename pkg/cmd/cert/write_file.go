package cert

import (
	"fmt"
	"os"
)

// writeFile writes to a file. If the file exists the operation will fail
func writeFile(filename string, data []byte) error {
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		return fmt.Errorf("file %s already exists", filename)
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err != nil {
		return err
	}
	defer f.Close()
	if n != len(data) {
		return fmt.Errorf("wrote %d of %d bytes", n, len(data))
	}
	return nil
}
