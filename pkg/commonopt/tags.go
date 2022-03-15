package commonopt

import "regexp"

// tagRegex matches tags of the form:
//
//   foo:bar
//   foo:"bar baz"
//   foo:
//   foo:""
//
var tagRegex = regexp.MustCompile(`^\s*(\S+):("?)(.*?)("?)\s*$`)

type Tags struct {
	Tags []string `long:"tag" description:"Set tag value (name:value)"`
}

func (t *Tags) AsMap() *map[string]string {
	tags := &map[string]string{}
	for _, elt := range t.Tags {
		res := tagRegex.FindStringSubmatch(elt)
		if len(res) != 5 {
			continue
		}
		(*tags)[res[1]] = res[3]
	}
	return tags
}
