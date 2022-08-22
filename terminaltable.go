package formatdata

import (
	"io"

	"github.com/shibukawa/stringwidth"
)

var tableVLine = []rune("│")[0]

const (
	tableTop    = 0
	tableHeader = 1
	tableMiddle = 2
	tableBottom = 3
)

var tableHorizonChars = [][]rune{
	[]rune("┌┬┐─"),
	[]rune("╞╪╡═"),
	[]rune("├┼┤─"),
	[]rune("└┴┘─"),
}

func renderSliceAsTerminalTable(table [][]any, tr *tableRenderer, eastAsianAmbiguousAsWide bool, out io.Writer) {
	// todo: EastAsianAmbiguous: true
	maxWidths, renderCells := calcTableSize(table, tr, eastAsianAmbiguousAsWide)
	repeat := func(r rune, length int) {
		for i := 0; i < length; i++ {
			io.WriteString(out, string(r))
		}
	}
	drawHorizontal := func(t int) {
		chars := tableHorizonChars[t]
		io.WriteString(out, tr.borderStart)
		io.WriteString(out, string(chars[0]))
		for i, m := range maxWidths {
			if i != 0 {
				io.WriteString(out, string(chars[1]))
			}
			repeat(chars[3], m+2)
		}
		io.WriteString(out, string(chars[2]))
		io.WriteString(out, tr.borderEnd)
		out.Write([]byte{'\n'})
	}

	drawHorizontal(tableTop)

	for r, row := range renderCells {
		io.WriteString(out, tr.border(string(tableVLine)))
		for i, w := range maxWidths {
			var c string
			if i < len(row) {
				c = row[i]
			}
			out.Write([]byte{' '})
			io.WriteString(out, c)
			repeat(' ', 1+w-stringwidth.Calc(c, stringwidth.Opt{
				IsAmbiguousWide: eastAsianAmbiguousAsWide,
			}))
			io.WriteString(out, tr.border(string(tableVLine)))
		}
		out.Write([]byte{'\n'})
		switch r {
		case len(renderCells) - 1:
			drawHorizontal(tableBottom)
		case 0:
			drawHorizontal(tableHeader)
		default:
			drawHorizontal(tableMiddle)
		}
	}
}
