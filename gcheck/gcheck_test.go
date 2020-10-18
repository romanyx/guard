package check

import (
	"fmt"
	"strings"
	"testing"
)

func TestMustNotNil(t *testing.T) {
	var nilVal interface{}
	var x map[string]string
	nilVal = x
	var ok interface{}
	ok = 1

	type args struct {
		nn int
		n  string
		i  interface{}
	}
	tt := []struct {
		name   string
		args   args
		expect string
	}{
		{
			name: "ok",
			args: args{1, "n", ok},
		},
		{
			name:   "nil",
			args:   args{1, "n", nil},
			expect: "parameter n number 1 is nil",
		},
		{
			name:   "nil value",
			args:   args{1, "n", nilVal},
			expect: "parameter n number 1 value is nil",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if r := recover(); r != nil {
					assert(t, tc.expect, fmt.Sprint(r))
				}
			}()
			MustNotNil(tc.args.nn, tc.args.n, tc.args.i)
		})
	}
}

func assert(t *testing.T, expect, got string) {
	t.Helper()

	if strings.Compare(expect, got) != 0 {
		t.Error("expected:")
		t.Errorf("%v", expect)
		t.Error("got:")
		t.Errorf("%v", got)
	}
}
