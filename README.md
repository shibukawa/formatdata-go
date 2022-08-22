# github.com/shibukawa/formatdata-go

[![Go Reference](https://pkg.go.dev/badge/github.com/shibukawa/formatdata-go.svg)](https://pkg.go.dev/github.com/shibukawa/formatdata-go)

Simple pretty print library

## Sample

```go
type SampleStruct struct {
	A string `json:"a"`
	B int    `json:"b"`
}

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
```

## API

`func FormatDataTo(data any, out io.Writer, o ...Opt) error`

There are three variations. ``o`` is option (described below)

* `func FormatData(data any, o ...Opt) error`

   This is an alias of `FormatDataTo(data, os.Stdout, o...)`

* `func FormatDataWithColor(data any, out io.Writer, o ...Opt) error`

   This always uses escape sequence to dump colorized output.

* `func FormatDataWithoutColor(data any, out io.Writer, o ...Opt) error`

   This always doesn't use escape sequence to dump colorized output.

### Option

```go
type Opt struct {
    // Terminal(default), Markdown, JSON, YAML
	OutputFormat             OutputFormat
    // Treat EastAsianAmbiguous characters as wide or not
	EastAsianAmbiguousAsWide bool
    // "terminal", "terminal8", "terminal16", "terminal256", "terminal16m". Default: "terminal"
	Formatter                string
    // https://github.com/alecthomas/chroma/tree/master/styles. Default: "monokai"
	Style                    string
    // Indent for JSON/YAML
	Indent                   int
}
```

If you want to specify YAML output, use this struct like this:

```go
formatdata.FormatData(d, formatdata.Opt{
    OutputFormat: formatdata.YAML
    Indent:       4,
})
```

## License

Apache 2

It uses code fragments of [github.com/alecthomas/chroma](https://github.com/alecthomas/chroma). This is developed by Alec Thomas under MIT license.
