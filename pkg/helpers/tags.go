package helpers

import "regexp"

// tagRegex matches tags of the form:
//
//   foo:bar
//   foo:"bar baz"
//   foo:
//   foo:""
//
var tagRegex = regexp.MustCompile(`^\s*(\S+):("?)(.*?)("?)\s*$`)

// TagMerge merges existing metadata with the tag array we get from the command
// line. If the tags parameter is nil a new map is created
func TagMerge(tags *map[string]string, args []string) *map[string]string {
	if tags == nil {
		tags = &map[string]string{}
	}
	for _, elt := range args {
		res := tagRegex.FindStringSubmatch(elt)
		if len(res) != 5 {
			continue
		}
		(*tags)[res[1]] = res[3]
	}
	return tags
}
