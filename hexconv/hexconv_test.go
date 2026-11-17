// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

package hexconv

import (
	"strings"
	"testing"
)

func charset(t *testing.T, name string) Charset {
	t.Helper()
	cs, ok := Lookup(name)
	if !ok {
		t.Fatalf("no charset %s", name)
	}
	return cs
}

func TestToHex(t *testing.T) {
	tests := []struct{ text, charset, want string }{
		{"Hello", "UTF-8", "48656c6c6f"},
		{"ğ", "UTF-8", "c49f"}, // two bytes in UTF-8
		{"ğ", "ISO-8859-9", "f0"},
		{"ğ", "windows-1254", "f0"},
		{"é", "ISO-8859-1", "e9"},
		{"A", "UTF-16", "feff0041"}, // byte order mark, big endian
		{"Hi", "US-ASCII", "4869"},
		{"", "UTF-8", ""},
	}
	for _, tt := range tests {
		got, err := ToHex(tt.text, charset(t, tt.charset))
		if err != nil || got != tt.want {
			t.Errorf("ToHex(%q, %s) = %q, %v; want %q", tt.text, tt.charset, got, err, tt.want)
		}
	}
}

func TestFromHex(t *testing.T) {
	tests := []struct{ hex, charset, want string }{
		{"f0", "ISO-8859-9", "ğ"},
		{"feff0041", "UTF-16", "A"},
		{"fffe4100", "UTF-16", "A"}, // little endian byte order mark
		{"0041", "UTF-16", "A"},     // big endian without a mark
		{"ff", "US-ASCII", "�"},
		{"ff", "UTF-8", "�"},
	}
	for _, tt := range tests {
		got, err := FromHex(tt.hex, charset(t, tt.charset))
		if err != nil || got != tt.want {
			t.Errorf("FromHex(%q, %s) = %q, %v; want %q", tt.hex, tt.charset, got, err, tt.want)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	for _, cs := range Charsets {
		texts := []string{"Hello", "", "line 1\nline 2"}
		switch cs.Name {
		case "UTF-8", "UTF-16":
			texts = append(texts, "Şişli'de ılık bir gün", "emoji 🙂")
		case "ISO-8859-9", "windows-1254":
			texts = append(texts, "Şişli'de ılık bir gün")
		}
		for _, text := range texts {
			h, err := ToHex(text, cs)
			if err != nil {
				t.Fatalf("%s: %v", cs, err)
			}
			back, err := FromHex(h, cs)
			if err != nil || back != text {
				t.Errorf("%s: %q -> %s -> %q, %v", cs, text, h, back, err)
			}
		}
	}
}

func TestSeparatorsAndUpperCase(t *testing.T) {
	for _, in := range []string{"48656C6C6F", "48 65 6c 6c 6f", "48:65:6c:6c:6f", "0x48 0x65 0x6C 0x6C 0x6F",
		"48-65-6c\n6c-6f", "48,65,6c,6c,6f", "\t48656c6c6f\r\n"} {
		if got, err := FromHex(in, UTF8); err != nil || got != "Hello" {
			t.Errorf("FromHex(%q) = %q, %v", in, got, err)
		}
	}
}

func TestOddNumberOfDigitsIsRejected(t *testing.T) {
	// The first version silently dropped the last digit and returned "He".
	_, err := FromHex("48656", UTF8)
	if err == nil || !strings.Contains(err.Error(), "odd number of digits (5)") {
		t.Errorf("err = %v", err)
	}
}

func TestNonHexCharactersAreRejected(t *testing.T) {
	// The first version silently removed them, so "hello" was read as the single digit "e".
	for _, in := range []string{"hello", "48zz", "48 6g"} {
		if _, err := FromHex(in, UTF8); err == nil || !strings.Contains(err.Error(), "not a hexadecimal digit") {
			t.Errorf("FromHex(%q): err = %v", in, err)
		}
	}
}

func TestEmptyInputGivesEmptyOutput(t *testing.T) {
	if got, err := FromHex("  \n ", UTF8); err != nil || got != "" {
		t.Errorf("got %q, %v", got, err)
	}
}

func TestCharactersOutsideTheCharsetAreReported(t *testing.T) {
	// Java wrote "?" for them, which turned into a different text on the way back.
	_, err := ToHex("Ağaç", charset(t, "US-ASCII"))
	if err == nil || !strings.Contains(err.Error(), `'ğ' (character 2) cannot be written in US-ASCII`) {
		t.Errorf("err = %v", err)
	}
	if _, err := ToHex("€", charset(t, "ISO-8859-9")); err == nil {
		t.Error("€ is not in ISO-8859-9")
	}
	if _, err := ToHex("€", charset(t, "windows-1254")); err != nil {
		t.Errorf("€ is in windows-1254: %v", err)
	}
}

func TestLookupIgnoresCase(t *testing.T) {
	if cs, ok := Lookup("utf-16"); !ok || cs.Name != "UTF-16" {
		t.Errorf("Lookup(utf-16) = %v, %v", cs, ok)
	}
	if _, ok := Lookup("EBCDIC"); ok {
		t.Error("EBCDIC is not supported")
	}
	if names := Names(); len(names) != 6 || names[0] != "UTF-8" {
		t.Errorf("Names() = %v", names)
	}
}
