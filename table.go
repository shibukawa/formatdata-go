package formatdata

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"

	"github.com/alecthomas/chroma/v2"
	"github.com/shibukawa/stringwidth"
	"gopkg.in/yaml.v3"
)

type tableRenderer struct {
	intCell     func(a int64, title bool) string
	uintCell    func(a uint64, title bool) string
	floatCell   func(a float64, title bool) string
	stringCell  func(a string, title bool) string
	boolCell    func(v bool, title bool) string
	otherCell   func(v any, title bool) string
	border      func(s string) string
	borderStart string
	borderEnd   string
}

func renderTable(cr *tableRenderer, cells [][]any, o Opt, out io.Writer) {
	if o.OutputFormat == Terminal {
		renderSliceAsTerminalTable(cells, cr, o.EastAsianAmbiguousAsWide, out)
	} else if o.OutputFormat == Markdown {
		renderSliceAsMarkdownTable(cells, cr, o.EastAsianAmbiguousAsWide, out)
	}
}

func newColorTextRenderer(style, formatter string) *tableRenderer {
	s := getStyle(style, formatter)

	findCategory := func(t chroma.TokenType) string {
		clr, ok := s[t]
		if !ok {
			clr, ok = s[t.SubCategory()]
			if !ok {
				clr = s[t.Category()]
			}
		}
		return clr
	}

	wrap := func(t chroma.TokenType, text string) string {
		clr := findCategory(t)
		if clr != "" {
			return clr + text + "\033[0m"
		}
		return text
	}
	// todo color
	return &tableRenderer{
		intCell: func(a int64, title bool) string {
			t := chroma.LiteralNumber
			if title {
				t = chroma.KeywordNamespace
			}
			return wrap(t, strconv.FormatInt(a, 10))
		},
		uintCell: func(a uint64, title bool) string {
			t := chroma.LiteralNumber
			if title {
				t = chroma.KeywordNamespace
			}
			return wrap(t, strconv.FormatUint(a, 10))
		},
		floatCell: func(a float64, title bool) string {
			t := chroma.LiteralNumber
			if title {
				t = chroma.KeywordNamespace
			}
			return wrap(t, strconv.FormatFloat(a, 'f', 6, 64))
		},
		stringCell: func(a string, title bool) string {
			t := chroma.LiteralNumber
			if title {
				t = chroma.KeywordNamespace
			}
			return wrap(t, a)
		},
		boolCell: func(a bool, title bool) string {
			t := chroma.LiteralStringBoolean
			if title {
				t = chroma.KeywordNamespace
			}
			if a {
				return wrap(t, "true")
			} else {
				return wrap(t, "false")
			}
		},
		otherCell: func(a any, title bool) string {
			t := chroma.LiteralNumber
			if title {
				t = chroma.KeywordNamespace
			}
			return wrap(t, fmt.Sprintf("%f", a))
		},
		border: func(s string) string {
			return wrap(chroma.LineTable, s)
		},
		borderStart: findCategory(chroma.LineTable),
		borderEnd:   "\033[0m",
	}
}

func newPlainTextTableRenderer() *tableRenderer {
	return &tableRenderer{
		intCell: func(a int64, title bool) string {
			return strconv.FormatInt(a, 10)
		},
		uintCell: func(a uint64, title bool) string {
			return strconv.FormatUint(a, 10)
		},
		floatCell: func(a float64, title bool) string {
			return strconv.FormatFloat(a, 'f', 6, 64)
		},
		stringCell: func(a string, title bool) string {
			return a
		},
		boolCell: func(a bool, title bool) string {
			if a {
				return "true"
			} else {
				return "false"
			}
		},
		otherCell: func(a any, title bool) string {
			return fmt.Sprintf("%f", a)
		},
		border: func(s string) string {
			return s
		},
		borderStart: "",
		borderEnd:   "",
	}
}

func calcTableSize(table [][]any, tr *tableRenderer, eastAsianAmbiguousAsWide bool) ([]int, [][]string) {
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	var maxWidths []int
	var renderCells [][]string
	for rowIndex, row := range table {
		renderRow := make([]string, len(row))
		for i, c := range row {
			switch v := c.(type) {
			case string:
				renderRow[i] = tr.stringCell(v, rowIndex == 0)
			case bool:
				renderRow[i] = tr.boolCell(v, rowIndex == 0)
			case int:
				renderRow[i] = tr.intCell(int64(v), rowIndex == 0)
			case int8:
				renderRow[i] = tr.intCell(int64(v), rowIndex == 0)
			case int16:
				renderRow[i] = tr.intCell(int64(v), rowIndex == 0)
			case int32:
				renderRow[i] = tr.intCell(int64(v), rowIndex == 0)
			case int64:
				renderRow[i] = tr.intCell(v, rowIndex == 0)
			case uint:
				renderRow[i] = tr.uintCell(uint64(v), rowIndex == 0)
			case uint8:
				renderRow[i] = tr.uintCell(uint64(v), rowIndex == 0)
			case uint16:
				renderRow[i] = tr.uintCell(uint64(v), rowIndex == 0)
			case uint32:
				renderRow[i] = tr.uintCell(uint64(v), rowIndex == 0)
			case uint64:
				renderRow[i] = tr.uintCell(v, rowIndex == 0)
			case float64:
				renderRow[i] = tr.floatCell(v, rowIndex == 0)
			case float32:
				renderRow[i] = tr.floatCell(float64(v), rowIndex == 0)
			default:
				renderRow[i] = tr.otherCell(v, rowIndex == 0)
			}
			if len(maxWidths) <= i {
				maxWidths = append(maxWidths, 0)
			}
			maxWidths[i] = max(maxWidths[i], stringwidth.Calc(renderRow[i], stringwidth.Opt{
				IsAmbiguousWide: eastAsianAmbiguousAsWide,
			}))
		}
		renderCells = append(renderCells, renderRow)
	}
	return maxWidths, renderCells
}

func canBeTable(data any) (cells [][]any, ok bool) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		cells, ok = convertToSliceOfSliceOfAny(v)
		if ok {
			return cells, ok
		}
	}

	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, false
	}
	var mapSlice []map[string]any
	err = yaml.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(&mapSlice)
	if err == nil {
		return canBeTable(mapSlice)
	}
	return nil, false
}

func isAllElement(slice reflect.Value, elementType reflect.Kind) bool {
	for i := 0; i < slice.Len(); i++ {
		if slice.Index(i).Kind() != elementType {
			return false
		}
	}
	return true
}

func convertToSliceOfSliceOfAny(slice reflect.Value) ([][]any, bool) {
	if isAllElement(slice, reflect.Slice) {
		var result [][]any
		for i := 0; i < slice.Len(); i++ {
			var row []any
			e := slice.Index(i)
			for j := 0; j < e.Len(); j++ {
				row = append(row, e.Index(j).Interface())
			}
			result = append(result, row)
		}
		return result, true
	} else if isAllElement(slice, reflect.Map) {
		var tempResult []map[string]any
		for i := 0; i < slice.Len(); i++ {
			e := slice.Index(i)
			tempRow := make(map[string]any)
			for _, k := range e.MapKeys() {
				if k.Kind() == reflect.String {
					tempRow[k.String()] = e.MapIndex(k).Interface()
				}
			}
			tempResult = append(tempResult, tempRow)
		}
		return convertTableMapToSlice(tempResult), true
	}
	return nil, false
}

func convertTableMapToSlice(table []map[string]any) [][]any {
	var headers []string
	existingCheck := map[string]bool{}
	for _, row := range table {
		for k := range row {
			if !existingCheck[k] {
				headers = append(headers, k)
				existingCheck[k] = true
			}
		}
	}
	sort.Strings(headers)

	slices := make([][]any, len(table)+1)
	headerSlice := make([]any, len(headers))
	for c, h := range headers {
		headerSlice[c] = h
	}
	slices[0] = headerSlice
	for r, row := range table {
		rowSlice := make([]any, len(headers))
		for c, h := range headers {
			if v, ok := row[h]; ok {
				rowSlice[c] = v
			} else {
				rowSlice[c] = ""
			}
		}
		slices[r+1] = rowSlice
	}
	return slices
}
