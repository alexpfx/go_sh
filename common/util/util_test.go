package util

import (
	"reflect"
	"testing"
)

func TestQuoteArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "t1",
			args: []string{
				"aabc aaabc",
			},
			want: []string{
				`"aabc aaabc"`,
			},
		},
		{
			name: "t2",
			args: []string{
				"aabc",
			},
			want: []string{
				"aabc",
			},
		},
		{
			name: "t3",
			args: []string{
				" aabc ",
			},
			want: []string{
				`" aabc "`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QuoteArgs(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuoteArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
