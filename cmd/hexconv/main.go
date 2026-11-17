// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

// Command hexconv converts on the command line, for scripts and terminals:
//
//	hexconv to-hex "Merhaba"
//	echo 4d 65 72 68 61 62 61 | hexconv to-text
//	hexconv -charset ISO-8859-9 to-hex ğ
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ergulcu/HexConverter/hexconv"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("hexconv", flag.ContinueOnError)
	flags.SetOutput(stderr)
	charset := flags.String("charset", "UTF-8", "character set: "+strings.Join(hexconv.Names(), ", "))
	flags.Usage = func() {
		fmt.Fprintln(stderr, "usage: hexconv [-charset NAME] to-hex|to-text [TEXT...]")
		fmt.Fprintln(stderr, "Without TEXT the input is read from stdin.")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() < 1 {
		flags.Usage()
		return 2
	}
	cs, ok := hexconv.Lookup(*charset)
	if !ok {
		fmt.Fprintf(stderr, "hexconv: unknown charset %q; use one of %s\n", *charset, strings.Join(hexconv.Names(), ", "))
		return 2
	}

	var convert func(string, hexconv.Charset) (string, error)
	switch flags.Arg(0) {
	case "to-hex":
		convert = hexconv.ToHex
	case "to-text":
		convert = hexconv.FromHex
	default:
		fmt.Fprintf(stderr, "hexconv: unknown command %q\n", flags.Arg(0))
		flags.Usage()
		return 2
	}

	var input string
	if flags.NArg() > 1 {
		input = strings.Join(flags.Args()[1:], " ")
	} else {
		b, err := io.ReadAll(stdin)
		if err != nil {
			fmt.Fprintln(stderr, "hexconv:", err)
			return 1
		}
		input = string(b)
		if flags.Arg(0) == "to-hex" {
			// "echo text | hexconv to-hex" should not encode the line break echo adds.
			input = strings.TrimSuffix(strings.TrimSuffix(input, "\n"), "\r")
		}
	}

	out, err := convert(input, cs)
	if err != nil {
		fmt.Fprintln(stderr, "hexconv:", err)
		return 1
	}
	fmt.Fprintln(stdout, out)
	return 0
}
