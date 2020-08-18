package main

import (
	"reflect"
	"testing"
)

func Test_uniqMsgs(t *testing.T) {
	type args struct {
		allMsgs []string
	}
	tests := []struct {
		name string
		args args
		want []Message
	}{
		{
			args: args{
				allMsgs: []string{"A", "B", "B"},
			},
			want: []Message{
				{Text: "A", Count: 1},
				{Text: "B", Count: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uniqMsgs(tt.args.allMsgs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uniqMsgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
