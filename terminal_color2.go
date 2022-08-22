package formatdata

import (
	"fmt"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
)

type style struct {
	Line string
}

var formatterColors = map[string]int{
	"terminal":    8,
	"terminal8":   8,
	"terminal16":  16,
	"terminal256": 256,
}

func getStyle(name, formatter string) map[chroma.TokenType]string {
	s := styles.Get(name)
	if formatter == "terminal16m" {
		return trueColorEscapeSequence(s)
	} else {
		colors, ok := formatterColors[formatter]
		if !ok {
			colors = 8
		}
		return styleToEscapeSequence(ttyTables[colors], s)
	}
}

func trueColorEscapeSequence(style *chroma.Style) map[chroma.TokenType]string {
	style = clearBackground(style)
	result := map[chroma.TokenType]string{}
	for _, ttype := range style.Types() {
		entry := style.Get(ttype)
		if !entry.IsZero() {
			out := ""
			if entry.Bold == chroma.Yes {
				out += "\033[1m"
			}
			if entry.Underline == chroma.Yes {
				out += "\033[4m"
			}
			if entry.Italic == chroma.Yes {
				out += "\033[3m"
			}
			if entry.Colour.IsSet() {
				out += fmt.Sprintf("\033[38;2;%d;%d;%dm", entry.Colour.Red(), entry.Colour.Green(), entry.Colour.Blue())
			}
			if entry.Background.IsSet() {
				out += fmt.Sprintf("\033[48;2;%d;%d;%dm", entry.Background.Red(), entry.Background.Green(), entry.Background.Blue())
			}
			result[ttype] = out
		}
	}
	return result

}
