# HexConverter

A JavaFX desktop tool that turns text into hexadecimal and back, in the character set you choose
(UTF-8, UTF-16, ISO-8859-1, ISO-8859-9, windows-1254, US-ASCII).

## Run

Java 25 is enough; JavaFX 27 is downloaded by Maven.

```bash
./mvnw javafx:run
```

A self-contained build with its own Java runtime:

```bash
./mvnw javafx:jlink
target/hexconverter/bin/hexconverter      # hexconverter.bat on Windows
```

## Hex input

Spaces, line breaks, `:`, `,`, `-` and `0x` prefixes between the bytes are accepted, upper or lower case:
`48 65 6c 6c 6f`, `48:65:6C:6C:6F` and `0x48,0x65` all work. Input with an odd number of digits or with
characters that are not hex digits is rejected with a message. The old version silently dropped the last
digit or the invalid characters and showed a wrong result.

## Tests

```bash
./mvnw verify
```

## License

Copyright 2019 Mehmet ERGULCU. Licensed under the [Apache License 2.0](LICENSE). Earlier versions were
published under the GPL v3. The Roboto Mono font in `src/main/resources/font` is by Google, also under
the Apache License 2.0.
