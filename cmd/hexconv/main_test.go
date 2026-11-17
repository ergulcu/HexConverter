// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"strings"
	"testing"
)

func runCLI(stdin string, args ...string) (code int, stdout, stderr string) {
	var out, errOut bytes.Buffer
	code = run(args, strings.NewReader(stdin), &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestArguments(t *testing.T) {
	if code, out, _ := runCLI("", "to-hex", "Hello", "World"); code != 0 || out != "48656c6c6f20576f726c64\n" {
		t.Errorf("code %d, out %q", code, out)
	}
	if code, out, _ := runCLI("", "-charset", "iso-8859-9", "to-text", "f0"); code != 0 || out != "ğ\n" {
		t.Errorf("code %d, out %q", code, out)
	}
}

func TestStdin(t *testing.T) {
	if code, out, _ := runCLI("Hello\n", "to-hex"); code != 0 || out != "48656c6c6f\n" {
		t.Errorf("the line break of echo must not be encoded: code %d, out %q", code, out)
	}
	if code, out, _ := runCLI("48 65\n6c 6c 6f\n", "to-text"); code != 0 || out != "Hello\n" {
		t.Errorf("code %d, out %q", code, out)
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		args []string
		code int
		msg  string
	}{
		{[]string{"to-text", "486"}, 1, "odd number of digits"},
		{[]string{"-charset", "US-ASCII", "to-hex", "ğ"}, 1, "cannot be written in US-ASCII"},
		{[]string{"-charset", "EBCDIC", "to-hex", "a"}, 2, "unknown charset"},
		{[]string{"reverse", "a"}, 2, "unknown command"},
		{nil, 2, "usage"},
	}
	for _, tt := range tests {
		code, _, errOut := runCLI("", tt.args...)
		if code != tt.code || !strings.Contains(errOut, tt.msg) {
			t.Errorf("%v: code %d, stderr %q", tt.args, code, errOut)
		}
	}
}
