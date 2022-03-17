package helpers

import (
	"encoding/json"

	"github.com/jedib0t/go-pretty/v6/table"
)

// DumpToTable writes a struct to a table writer with one row for each field. Arrays are not (yet)
// handled.
func DumpToTable(t table.Writer, s any) error {
	// Dump as name.value list by converting to JSON and then back to map[string]any
	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}

	nameValue := map[string]any{}
	if err := json.Unmarshal(buf, &nameValue); err != nil {
		return err
	}
	dumpFields(t, "", nameValue)
	return nil
}

func fieldName(prefix string, field string) string {
	if prefix == "" {
		return field
	}
	return prefix + "." + field
}

func dumpFields(t table.Writer, prefix string, nameValue map[string]any) {
	for k, v := range nameValue {
		subType, ok := v.(map[string]any)
		if ok {
			dumpFields(t, fieldName(prefix, k), subType)
			continue
		}
		t.AppendRow(table.Row{
			fieldName(prefix, k),
			v,
		})
	}
}
