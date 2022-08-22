package formatdata

import (
	"io"

	"github.com/shibukawa/stringwidth"
)

func renderSliceAsMarkdownTable(table [][]any, cr *tableRenderer, eastAsianAmbiguousAsWide bool, out io.Writer) {
	maxWidths, renderCells := calcTableSize(table, cr, eastAsianAmbiguousAsWide)
	repeat := func(r rune, length int) {
		for i := 0; i < length; i++ {
			io.WriteString(out, string(r))
		}
	}
	drawHorizontal := func() {
		out.Write([]byte{'|'})
		for i, m := range maxWidths {
			if i != 0 {
				out.Write([]byte{'|'})
			}
			repeat('-', m+2)
		}
		out.Write([]byte{'|', '\n'})
	}

	for r, row := range renderCells {
		out.Write([]byte{'|'})
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
			out.Write([]byte{'|'})
		}
		out.Write([]byte{'\n'})
		if r == 0 && len(renderCells) != 1 {
			drawHorizontal()
		}
	}
}
