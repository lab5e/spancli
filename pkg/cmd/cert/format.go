package cert

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// formatBytes formats the byte buffer into a suitable format. The input is PEM-encoded data
func formatBytes(rawBytes []byte, format string) ([]byte, error) {
	var derBytes []byte
	block, remain := pem.Decode(rawBytes)
	for block != nil {
		derBytes = append(derBytes, block.Bytes...)
		block, remain = pem.Decode(remain)
	}

	var outBytes []byte

	switch format {

	case "der":
		outBytes = derBytes

	case "base64":
		outBytes = []byte(base64.StdEncoding.EncodeToString(derBytes))

	case "hex-list":
		for i, v := range derBytes {
			if i == 0 {
				outBytes = append(outBytes, []byte(fmt.Sprintf("0x%02x", v))...)
				continue
			}
			outBytes = append(outBytes, []byte(fmt.Sprintf(",0x%02x", v))...)
		}

	case "decimal-list":
		for i, v := range derBytes {
			if i == 0 {
				outBytes = append(outBytes, []byte(fmt.Sprintf("%d", v))...)
				continue
			}
			outBytes = append(outBytes, []byte(fmt.Sprintf(",%d", v))...)
		}
	default:
		return nil, fmt.Errorf("unknown format string: %s", format)
	}

	return outBytes, nil
}
