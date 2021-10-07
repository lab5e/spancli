package main

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// SpanCLIStyle is a custom style for the mesh CLI
var SpanCLIStyle = table.Style{
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

// newTableOutput creates a new table writer with the specified settings
func newTableOutput(format string, no_color bool, pagesize int) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	if !no_color {
		t.SetStyle(SpanCLIStyle)
	}

	if pagesize != 0 {
		t.SetPageSize(pagesize)
	}

	return t
}

// renderTable renders the table according to settings
func renderTable(t table.Writer, format string) {
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
