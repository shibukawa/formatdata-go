package formatdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trimIndent(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "leading newlines",
			args: args{
				src: `
                     test1
                     test2`,
			},
			want: `test1
test2`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, trimIndent(tt.args.src), "trimIndent(%v)", tt.args.src)
		})
	}
}
