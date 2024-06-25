package functions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunFuncWithArgs(t *testing.T) {
	output := []string{}
	type args struct {
		funcsAndArgs []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "TestRunFuncWithArgs",
			args: args{
				funcsAndArgs: []interface{}{
					func() error {
						output = append(output, "test1")
						return nil
					},
					func(in, out string) error {
						output = append(output, in, out)
						return nil
					},
					"test2", "test3",
					func() error {
						output = append(output, "test4")
						return nil
					},
				},
			},
			want:    []string{"test1", "test2", "test3", "test4"},
			wantErr: false,
		},
		{
			name: "TestRunFuncWithArgsWithErrors",
			args: args{
				funcsAndArgs: []interface{}{
					func() error {
						output = append(output, "test1")
						return nil
					},
					func(in, out string) error {
						output = append(output, in, out)
						return nil
					},
					"test2",
					func() error {
						output = append(output, "test4")
						return nil
					},
				},
			},
			want:    []string{"test1"},
			wantErr: true,
		},
		{
			name: "TestRunFuncWithArgsWithVariadicArgs",
			args: args{
				funcsAndArgs: []interface{}{
					func() error {
						output = append(output, "test1")
						return nil
					},
					func(in, out string) error {
						output = append(output, in, out)
						return nil
					},
					"test2", "test3",
					func(a ...string) error {
						fmt.Println(a)
						output = append(output, a...)
						return nil
					},
					"test4",
					"test5",
					"test6",
				},
			},
			want:    []string{"test1", "test2", "test3", "test4", "test5", "test6"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunFuncWithArgs(tt.args.funcsAndArgs...); (err != nil) != tt.wantErr {
				t.Errorf("RunFuncWithArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, output)
			output = []string{}
		})
	}
}
