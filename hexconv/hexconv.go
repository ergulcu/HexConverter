// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

// Package hexconv turns text into hexadecimal in a chosen character set and back.
package hexconv

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Charset is a character set text can be written in.
type Charset struct {
	Name string
	enc  encoding.Encoding
}

func (c Charset) String() string { return c.Name }

// UTF8 is the default character set.
var UTF8 = Charset{"UTF-8", unicode.UTF8}

// Charsets lists the supported character sets, UTF-8 first. ISO-8859-9 and windows-1254 cover Turkish text
// written by older systems. UTF-16 writes a byte order mark and reads big endian when there is none.
var Charsets = []Charset{
	UTF8,
	{"UTF-16", unicode.UTF16(unicode.BigEndian, unicode.UseBOM)},
	{"ISO-8859-1", charmap.ISO8859_1},
	{"ISO-8859-9", charmap.ISO8859_9},
	{"windows-1254", charmap.Windows1254},
	{"US-ASCII", ascii{}},
}

// Names returns the names of the supported character sets.
func Names() []string {
	names := make([]string, len(Charsets))
	for i, c := range Charsets {
		names[i] = c.Name
	}
	return names
}

// Lookup finds a character set by name, ignoring case.
func Lookup(name string) (Charset, bool) {
	for _, c := range Charsets {
		if strings.EqualFold(c.Name, name) {
			return c, true
		}
	}
	return Charset{}, false
}

// ToHex encodes the text in the character set and returns the bytes as lower case hex, two digits per byte.
// It fails when the text has a character the character set cannot hold.
func ToHex(text string, cs Charset) (string, error) {
	b, err := cs.enc.NewEncoder().Bytes([]byte(text))
	if err != nil {
		return "", unmappable(text, cs)
	}
	return hex.EncodeToString(b), nil
}

// FromHex reads hex digits and decodes the bytes in the character set. Separators between the bytes are allowed:
// white space, ':', ',', '-' and a "0x" in front of a byte, so "48 65", "48:65" and "0x48,0x65" all work.
// Bytes that are not valid in the character set become U+FFFD.
func FromHex(s string, cs Charset) (string, error) {
	digits, err := Normalize(s)
	if err != nil {
		return "", err
	}
	b, err := hex.DecodeString(digits)
	if err != nil {
		return "", err // cannot happen after Normalize
	}
	text, err := cs.enc.NewDecoder().Bytes(b)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

var separators = regexp.MustCompile(`0[xX]|[\s:,\-]`)

// Normalize removes the separators and checks that only complete pairs of hex digits remain.
func Normalize(s string) (string, error) {
	digits := separators.ReplaceAllString(s, "")
	for _, r := range digits {
		if !isHexDigit(r) {
			return "", fmt.Errorf("%q is not a hexadecimal digit; only 0-9 and a-f are allowed", r)
		}
	}
	if len(digits)%2 != 0 {
		return "", fmt.Errorf("the hex text has an odd number of digits (%d); every byte needs two digits", len(digits))
	}
	return digits, nil
}

func isHexDigit(r rune) bool {
	return '0' <= r && r <= '9' || 'a' <= r && r <= 'f' || 'A' <= r && r <= 'F'
}

// unmappable names the first character the character set cannot hold.
func unmappable(text string, cs Charset) error {
	enc := cs.enc.NewEncoder()
	for i, r := range text {
		if _, err := enc.String(string(r)); err != nil {
			return fmt.Errorf("%q (character %d) cannot be written in %s", r, utf8.RuneCountInString(text[:i])+1, cs.Name)
		}
	}
	return fmt.Errorf("the text cannot be written in %s", cs.Name)
}

// ascii is US-ASCII: characters above 127 cannot be encoded and such bytes decode to U+FFFD.
type ascii struct{}

var errNotASCII = errors.New("not US-ASCII")

func (ascii) NewEncoder() *encoding.Encoder {
	return &encoding.Encoder{Transformer: asciiTransformer{encode: true}}
}

func (ascii) NewDecoder() *encoding.Decoder {
	return &encoding.Decoder{Transformer: asciiTransformer{}}
}

type asciiTransformer struct{ encode bool }

func (t asciiTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for nSrc < len(src) {
		c := src[nSrc]
		if c < utf8.RuneSelf {
			if nDst >= len(dst) {
				return nDst, nSrc, transform.ErrShortDst
			}
			dst[nDst] = c
			nDst++
			nSrc++
			continue
		}
		if t.encode {
			return nDst, nSrc, errNotASCII
		}
		if nDst+3 > len(dst) {
			return nDst, nSrc, transform.ErrShortDst
		}
		nDst += copy(dst[nDst:], "�")
		nSrc++
	}
	return nDst, nSrc, nil
}

func (asciiTransformer) Reset() {}
