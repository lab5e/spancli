package commonopt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var parms = []string{
	``,
	`  `,
	`nonsense`,
	`nonsense `,
	` nonsense`,
	`foo:bar`,
	`foo:"bar"`,
	`foo:"bar baz"`,
	`fooBar:"bar"`,
	`foo-bar:baz`,
	`foo_bar:baz`,
	`existing:`,
	`otherexisting:""`,
}

var meta = &map[string]string{
	"foo":           "something else",
	"existing":      "some value",
	"otherexisting": "some other value",
	"untouched":     "still here",
}

func TestTagRegex(t *testing.T) {
	for _, s := range parms {
		res := tagRegex.FindStringSubmatch(s)
		if len(res) == 0 {
			assert.True(t, strings.TrimSpace(s) == "" || strings.TrimSpace(s) == "nonsense")
			continue
		}

		assert.NotZero(t, len(res[1]))
		assert.True(t, len(res[3]) >= 0)
	}
}

func TestMetaToMap(t *testing.T) {
	tt := Tags{Tags: parms}
	mm := tt.AsMap()

	assert.Len(t, *mm, 6)
	assert.Equal(t, (*mm)["foo"], "bar baz")
	assert.Equal(t, (*mm)["fooBar"], "bar")
	assert.Equal(t, (*mm)["foo-bar"], "baz")
	assert.Equal(t, (*mm)["foo_bar"], "baz")
}
