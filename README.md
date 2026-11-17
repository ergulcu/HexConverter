# HexConverter

Turns text into hexadecimal and back, in the character set you choose (UTF-8, UTF-16, ISO-8859-1,
ISO-8859-9, windows-1254, US-ASCII). Written in Go: the window of the earlier JavaFX version (kept at the
[`java` tag](https://github.com/ergulcu/HexConverter/tree/java)), plus a command line tool, both as single
binaries without a Java runtime.

## Desktop program

The window uses [Fyne](https://fyne.io), which needs cgo: a C compiler (gcc or clang), and on Linux the OpenGL,
X11 and Wayland headers (`sudo apt-get install gcc libgl1-mesa-dev xorg-dev libxkbcommon-dev libwayland-dev`). On Windows,
[MSYS2](https://www.msys2.org) or [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) provide gcc.

```bash
go run ./cmd/hexconverter
# a Windows build without a console window:
go build -ldflags="-H=windowsgui" -o hexconverter.exe ./cmd/hexconverter
```

Ready builds for Linux, Windows and macOS are attached to every CI run (Actions, the `hexconverter-*` artifacts).

Type text or hex on the left, choose the character set and press **Text To Hex** or **Hex To Text**; the
result appears on the right and **Copy Result** puts it on the clipboard.

## Command line

Pure Go, so it builds anywhere without a C compiler:

```bash
go install github.com/ergulcu/HexConverter/cmd/hexconv@latest

hexconv to-hex "Merhaba"                       # 4d657268616261
echo "4d 65 72 68 61 62 61" | hexconv to-text  # Merhaba
hexconv -charset ISO-8859-9 to-hex ğ           # f0
```

Without text after the command, the input is read from stdin (the line break `echo` adds is not encoded).
The exit code is 1 when the input cannot be converted and 2 for wrong usage.

## Hex input

Spaces, line breaks, `:`, `,`, `-` and `0x` prefixes between the bytes are accepted, upper or lower case:
`48 65 6c 6c 6f`, `48:65:6C:6C:6F` and `0x48,0x65` all work. Input with an odd number of digits or with
characters that are not hex digits is rejected with a message.

## Differences from the Java version

- A character the chosen character set cannot hold (`ğ` in US-ASCII) is reported. The Java version wrote
  `?` instead, so converting back gave a different text.
- UTF-16 works as in Java: a byte order mark is written, and input without one is read big endian.
- The converter is a Go package, [`hexconv`](hexconv), that other programs can import.

## Tests

```bash
go test ./hexconv ./internal/... ./cmd/hexconv   # the window tests run headless, without cgo
```

## License

Copyright 2019 Mehmet ERGULCU. Licensed under the [Apache License 2.0](LICENSE).
