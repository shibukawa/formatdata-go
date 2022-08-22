//go:build !windows

package formatdata

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// FormatDataTo is variation of [FormatData]. You can specify output destination.
//
// If out is terminal, it uses escape sequence to dump colorized output.
func FormatDataTo(data any, out io.Writer, o ...Opt) error {
	if fo, ok := out.(*os.File); ok {
		if terminal.IsTerminal(int(fo.Fd())) {
			return FormatDataWithColor(data, out, o...)
		}
	}
	return FormatDataWithoutColor(data, out, o...)
}
