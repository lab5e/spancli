package helpers

import (
	"fmt"
	"strings"
)

// TagsToString converts a map of tags into a string with key:value fields
func TagsToString(tags map[string]string) string {
	tagStrs := []string{}
	for k, v := range tags {
		tagStrs = append(tagStrs, fmt.Sprintf("%s:%s", k, v))
	}

	return strings.Join(tagStrs, ",")
}
