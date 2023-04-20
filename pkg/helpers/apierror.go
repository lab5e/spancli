package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// APIError creates an error instance based on error message and HTTP response returned from API call.
func APIError(res *http.Response, e error) error {
	var body []byte
	var err error

	if res != nil {
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			// If we can't read the body, just return the error from the API
			return e
		}
	}

	var errmsg struct {
		Status  int    `json:"status"`
		Message string `json:"title"`
	}

	err = json.Unmarshal(body, &errmsg)
	if err != nil {
		// If we can't extract the error message we just return the error from the API
		return e
	}

	return fmt.Errorf("%s: %s", res.Status, errmsg.Message)
}
