package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get pwd: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		expect      string
		expectError string
	}{
		{
			name:   "base",
			args:   []string{"func", "NewStr"},
			expect: "// guardNewStr allows to guard NewStr constructor.\nfunc guardNewStr(r io.Reader, o other, a *arg, m map[string]string, x string) {\n\tgcheck.MustNotNil(1, \"r\", r)\n\tgcheck.MustNotNil(2, \"o\", o)\n\tgcheck.MustNotNil(3, \"a\", a)\n\tgcheck.MustNotNil(4, \"m\", m)\n\n}\n",
		},
		{
			name:        "usage",
			args:        []string{},
			expectError: "go-guard [action] <function>\ngo-guard generates guard function for a constructor\nor constuctor call\nExamples:\ngo-guard func NewRepository\ngo-guard call NewRepository\nexit status 2\n",
		},
		{
			name:        "unknow func",
			args:        []string{"func", "unknown"},
			expectError: "function not found or has no parameters\nexit status 2\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			args := []string{
				"run",
				filepath.Join(strings.TrimSuffix(
					pwd,
					"/tests",
				), "main.go"),
			}
			args = append(args, tt.args...)
			cmd := exec.Command("go", args...)

			var buf bytes.Buffer
			cmd.Stdout = &buf
			cmd.Stderr = &buf

			if err := cmd.Run(); err != nil {
				assert(t, tt.expectError, buf.String())
				return
			}

			assert(t, tt.expect, buf.String())
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
