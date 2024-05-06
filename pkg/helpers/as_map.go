package helpers

import "regexp"

// configRegex matches tags of the form:
//
//	foo:bar
//	foo:"bar baz"
//	foo:
//	foo:""
var configRegex = regexp.MustCompile(`^\s*(\S+):("?)(.*?)("?)\s*$`)

// AsMap converts a list of strings with name:value (with optional quotes) into a
// *map[string]string
func AsMap(config []string) *map[string]string {
	params := make(map[string]string)

	for _, elt := range config {
		res := configRegex.FindStringSubmatch(elt)
		if len(res) != 5 {
			continue
		}
		params[res[1]] = res[3]
	}
	return &params
}
