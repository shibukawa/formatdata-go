package formatdata

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("🐴")
}

func TestFormatData(t *testing.T) {
	type args struct {
		data any
		opt  Opt
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "Terminal(default): table ok data",
			args: args{
				data: [][]string{
					{
						"AAAA",
						"BBB",
						"CCCC",
					},
					{
						"12",
						"123",
						"1234",
					},
				},
				opt: Opt{},
			},
			wantOut: trimIndent(`
				┌──────┬─────┬──────┐
				│ AAAA │ BBB │ CCCC │
				╞══════╪═════╪══════╡
				│ 12   │ 123 │ 1234 │
				└──────┴─────┴──────┘
				`),
		},
		{
			name: "Terminal(default): table ng data -> fallback to YAML",
			args: args{
				data: []string{
					"AAAA",
					"BBB",
					"CCCC",
				},
				opt: Opt{},
			},
			wantOut: trimIndent(`
                   - AAAA
                   - BBB
                   - CCCC
`),
		},
		{
			name: "Terminal: table ok data",
			args: args{
				data: [][]string{
					{
						"AAAA",
						"BBB",
						"CCCC",
					},
					{
						"12",
						"123",
						"1234",
					},
				},
				opt: Opt{
					OutputFormat: Terminal,
				},
			},
			wantOut: trimIndent(`
				┌──────┬─────┬──────┐
				│ AAAA │ BBB │ CCCC │
				╞══════╪═════╪══════╡
				│ 12   │ 123 │ 1234 │
				└──────┴─────┴──────┘
				`),
		},
		{
			name: "Terminal: table ng data -> fallback to YAML",
			args: args{
				data: []string{
					"AAAA",
					"BBB",
					"CCCC",
				},
				opt: Opt{
					OutputFormat: Terminal,
				},
			},
			wantOut: trimIndent(`
                   - AAAA
                   - BBB
                   - CCCC
`),
		},
		{
			name: "Markdown: table ok data",
			args: args{
				data: [][]string{
					{
						"AAAA",
						"BBB",
						"CCCC",
					},
					{
						"12",
						"123",
						"1234",
					},
				},
				opt: Opt{
					OutputFormat: Markdown,
				},
			},
			wantOut: trimIndent(`
				| AAAA | BBB | CCCC |
				|------|-----|------|
				| 12   | 123 | 1234 |
				`),
		},
		{
			name: "Markdown: table ng data -> fallback to YAML",
			args: args{
				data: []string{
					"AAAA",
					"BBB",
					"CCCC",
				},
				opt: Opt{
					OutputFormat: Markdown,
				},
			},
			wantOut: trimIndent(`
                   - AAAA
                   - BBB
                   - CCCC
`),
		},
		{
			name: "YAML: table ok data",
			args: args{
				data: [][]string{
					{
						"AAAA",
						"BBB",
						"CCCC",
					},
					{
						"12",
						"123",
						"1234",
					},
				},
				opt: Opt{
					OutputFormat: YAML,
				},
			},
			wantOut: trimIndent(`
				- - AAAA
				  - BBB
				  - CCCC
				- - "12"
				  - "123"
				  - "1234"
`),
		},
		{
			name: "JSON: table ok data",
			args: args{
				data: [][]string{
					{
						"AAAA",
						"BBB",
						"CCCC",
					},
					{
						"12",
						"123",
						"1234",
					},
				},
				opt: Opt{
					OutputFormat: JSON,
				},
			},
			wantOut: trimIndent(`
				[
				  [
				    "AAAA",
				    "BBB",
				    "CCCC"
				  ],
				  [
				    "12",
				    "123",
				    "1234"
				  ]
				]
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			FormatDataWithoutColor(tt.args.data, out, tt.args.opt)
			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
