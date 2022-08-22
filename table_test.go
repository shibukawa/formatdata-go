package formatdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type SampleStruct struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func Test_canBeTable(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name   string
		args   args
		want   [][]any
		wantOk bool
	}{
		{
			name: "primitive can't be table",
			args: args{
				data: 1,
			},
			wantOk: false,
		},
		{
			name: "nested array can be table",
			args: args{
				data: [][]int{{1, 2}, {3, 4}},
			},
			want:   [][]any{{1, 2}, {3, 4}},
			wantOk: true,
		},
		{
			name: "slice of map can be table",
			args: args{
				data: []map[string]int{{"a": 1, "b": 2}, {"a": 3, "b": 4}},
			},
			want:   [][]any{{"a", "b"}, {1, 2}, {3, 4}},
			wantOk: true,
		},
		{
			name: "slice of JSON compatible struct can be table",
			args: args{
				data: []SampleStruct{{A: "1", B: 2}, {A: "3", B: 4}},
			},
			want:   [][]any{{"a", "b"}, {"1", 2}, {"3", 4}},
			wantOk: true,
		},
		{
			name: "false if in other cases",
			args: args{
				data: map[string]string{
					"not nested": "return false",
				},
			},
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCells, gotOk := canBeTable(tt.args.data)
			assert.Equal(t, tt.want, gotCells)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}
