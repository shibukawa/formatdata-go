// package formatdata provides pretty-print function [FormatData] and its variations
package formatdata

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/mattn/go-colorable"
	"gopkg.in/yaml.v3"
)

type OutputFormat = int

const (
	Terminal OutputFormat = iota // Default. If data is not grid compatible, fallback to YAML.
	Markdown                     // Markdown table. If data is not grid compatible, fallback to YAML.
	JSON
	YAML
)

type Opt struct {
	OutputFormat             OutputFormat
	EastAsianAmbiguousAsWide bool
	Formatter                string // "terminal", "terminal8", "terminal16", "terminal256", "terminal16m". Default: "terminal"
	Style                    string // https://github.com/alecthomas/chroma/tree/master/styles. Default: "monokai"
	Indent                   int    // Indent for JSON/YAML
}

// FormatData is the simplest API.
//
// data is any data to show. Nested array or array of struct will be formatted in table.
//
// o is an optional. It controls format, style and so on.
func FormatData(data any, o ...Opt) error {
	return FormatDataTo(data, colorable.NewColorableStdout(), o...)
}

func normalizeOpt(o []Opt) Opt {
	var result Opt
	if len(o) > 0 {
		result = o[0]
	}
	if result.Indent == 0 {
		result.Indent = 2
	}
	if result.Formatter == "" {
		result.Formatter = "terminal"
	}
	if result.Style == "" {
		result.Style = "monokai"
	}
	return result
}

// FormatDataWithColor is [FormatDataTo]'s variation that always uses escape sequence to dump colorized output.
func FormatDataWithColor(data any, out io.Writer, o ...Opt) error {
	opt := normalizeOpt(o)
	if opt.OutputFormat == Terminal || opt.OutputFormat == Markdown {
		if cells, ok := canBeTable(data); ok {
			cr := newColorTextRenderer(opt.Style, opt.Formatter)
			if opt.OutputFormat == Terminal {
				renderSliceAsTerminalTable(cells, cr, opt.EastAsianAmbiguousAsWide, out)
			} else {
				renderSliceAsMarkdownTable(cells, cr, opt.EastAsianAmbiguousAsWide, out)
			}
			return nil
		} else {
			opt.OutputFormat = YAML
		}
	}
	if opt.OutputFormat == YAML {
		var b bytes.Buffer
		e := yaml.NewEncoder(&b)
		e.SetIndent(opt.Indent)
		e.Encode(data)
		return quick.Highlight(out, b.String(), "yaml", opt.Formatter, opt.Style)
	} else /* JSON */ {
		var b bytes.Buffer
		e := json.NewEncoder(&b)
		e.SetIndent("", strings.Repeat(" ", opt.Indent))
		return quick.Highlight(out, b.String(), "json", opt.Formatter, opt.Style)
	}
}

// FormatDataWithColor is [FormatDataTo]'s variation that always doesn't use escape sequence.
func FormatDataWithoutColor(data any, out io.Writer, o ...Opt) error {
	opt := normalizeOpt(o)
	if opt.OutputFormat == Terminal || opt.OutputFormat == Markdown {
		if cells, ok := canBeTable(data); ok {
			cr := newPlainTextTableRenderer()
			if opt.OutputFormat == Terminal {
				renderSliceAsTerminalTable(cells, cr, opt.EastAsianAmbiguousAsWide, out)
			} else {
				renderSliceAsMarkdownTable(cells, cr, opt.EastAsianAmbiguousAsWide, out)
			}
			return nil
		} else {
			opt.OutputFormat = YAML
		}
	}
	if opt.OutputFormat == YAML {
		e := yaml.NewEncoder(out)
		e.SetIndent(opt.Indent)
		return e.Encode(data)
	} else /* JSON */ {
		e := json.NewEncoder(out)
		e.SetIndent("", strings.Repeat(" ", opt.Indent))
		return e.Encode(data)
	}
}
