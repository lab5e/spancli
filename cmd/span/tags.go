package main

import "regexp"

// tagRegex matches tags of the form:
//
//   foo:bar
//   foo:"bar baz"
//   foo:
//   foo:""
//
var tagRegex = regexp.MustCompile(`^\s*(\S+):("?)(.*?)("?)\s*$`)

// tagMerge merges existing metadata with the tag array we get from the command
// line.
func tagMerge(tags *map[string]string, args []string) *map[string]string {
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

// tagsFromArgs creates a tag map from the arguments
func tagsFromArgs(args []string) *map[string]string {
	return tagMerge(nil, args)
}
