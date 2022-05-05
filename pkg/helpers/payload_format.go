package helpers

import "encoding/base64"

// PayloadFormat formats the payload
func PayloadFormat(pl string, decode bool) string {
	if decode {
		ret, err := base64.StdEncoding.DecodeString(pl)
		if err != nil {
			return "(error)"
		}
		return string(ret)
	}
	return pl
}
