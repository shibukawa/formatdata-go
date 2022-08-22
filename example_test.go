package formatdata_test

import (
	"github.com/shibukawa/formatdata-go"
)

type SampleStruct struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func ExampleFormatData_TableStyle() {
	// If the data can be represented as table, it shows in table format
	tableData := []SampleStruct{
		{
			A: "a",
			B: 1,
		},
		{
			A: "b",
			B: 2,
		},
	}
	formatdata.FormatData(tableData)
	// Output:
	// ┌───┬───┐
	// │ a │ b │
	// ╞═══╪═══╡
	// │ a │ 1 │
	// ├───┼───┤
	// │ b │ 2 │
	// └───┴───┘
}

func ExampleFormatData_MarkdownTable() {
	// If the data can be represented as table and Markdown is specified,
	// it shows in Markdown table format. Fallback format is YAML.
	tableData := []SampleStruct{
		{
			A: "a",
			B: 1,
		},
		{
			A: "b",
			B: 2,
		},
	}
	formatdata.FormatData(tableData, formatdata.Opt{
		OutputFormat: formatdata.Markdown,
	})
	// Output:
	// | a | b |
	// |---|---|
	// | a | 1 |
	// | b | 2 |
}

func ExampleFormatData_YAML() {
	// If the data can not be represented as table and no format specified, or YAML is specified,
	// it shows in YAML format
	tableData := []SampleStruct{
		{
			A: "a",
			B: 1,
		},
		{
			A: "b",
			B: 2,
		},
	}
	formatdata.FormatData(tableData, formatdata.Opt{
		OutputFormat: formatdata.YAML,
	})
	// Output:
	// - a: a
	//   b: 1
	// - a: b
	//   b: 2
}
