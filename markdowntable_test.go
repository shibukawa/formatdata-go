package formatdata

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownRenderer_RenderSliceAsTable(t1 *testing.T) {
	type args struct {
		wideFlag bool
		cells    [][]any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple string tables",
			args: args{
				wideFlag: false,
				cells: [][]any{
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
					{
						"1234",
						"123",
						"12",
					},
				},
			},
			want: trimIndent(`
				| AAAA | BBB | CCCC |
				|------|-----|------|
				| 12   | 123 | 1234 |
				| 1234 | 123 | 12   |
				`),
		},
		{
			name: "simple string tables",
			args: args{
				wideFlag: false,
				cells: [][]any{
					{
						"a",
						"b",
					},
					{
						"a",
						"1",
					},
					{
						"b",
						"2",
					},
				},
			},
			want: trimIndent(`
				| a | b |
				|---|---|
				| a | 1 |
				| b | 2 |
				`),
		},
		{
			name: "invalid column length tables",
			args: args{
				wideFlag: false,
				cells: [][]any{
					{
						"AAAA",
						"BBB",
						"CCCC",
						"DDDDD",
					},
					{
						"12",
						"123",
						"1234",
					},
					{
						"1234",
						"123",
						"12",
					},
				},
			},
			want: trimIndent(`
				| AAAA | BBB | CCCC | DDDDD |
				|------|-----|------|-------|
				| 12   | 123 | 1234 |       |
				| 1234 | 123 | 12   |       |
				`),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			renderSliceAsMarkdownTable(tt.args.cells, newPlainTextTableRenderer(), false, &buf)
			assert.Equalf(t1, tt.want, buf.String(), "renderSliceAsMarkdownTable(%v, ...)", tt.args.cells)
		})
	}
}
