package helpers

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// SpanCLIStyle is the table style
var spanCLIStyle = table.Style{
	Name: "SpanCLIStyle",
	Box:  table.StyleBoxLight,
	Color: table.ColorOptions{
		IndexColumn:  text.Colors{text.BgHiRed, text.FgBlack},
		Footer:       text.Colors{text.BgRed, text.FgBlack},
		Header:       text.Colors{text.BgRed, text.FgWhite},
		Row:          text.Colors{text.BgHiWhite, text.FgBlack},
		RowAlternate: text.Colors{text.BgWhite, text.FgBlack},
	},
	Format: table.FormatOptions{
		Footer: text.FormatDefault,
		Header: text.FormatDefault,
		Row:    text.FormatDefault,
	},

	HTML:    table.DefaultHTMLOptions,
	Options: table.OptionsNoBordersAndSeparators,
	Title: table.TitleOptions{
		Align:  text.AlignLeft,
		Colors: text.Colors{text.BgRed, text.FgBlack},
		Format: text.FormatDefault,
	},
}

// NewTableOutput creates a new table writer with the specified settings
func NewTableOutput(format string, no_color bool, pagesize int) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	if !no_color {
		t.SetStyle(spanCLIStyle)
	}

	if pagesize != 0 {
		t.SetPageSize(pagesize)
	}

	return t
}

// RenderTable renders the table according to settings
func RenderTable(t table.Writer, format string) {
	switch format {
	case "csv":
		t.SetTitle("")
		t.RenderCSV()
	case "html":
		t.RenderHTML()
	case "markdown":
		t.RenderMarkdown()
	default:
		fmt.Println()
		t.Render()
		fmt.Println()
	}
}
