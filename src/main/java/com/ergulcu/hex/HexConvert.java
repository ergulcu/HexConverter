// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

package com.ergulcu.hex;

import java.nio.charset.Charset;
import java.util.HexFormat;
import java.util.Objects;
import java.util.regex.Pattern;

/** Converts text to its hexadecimal form in a given charset and back. */
public final class HexConvert {

    /** Characters that only separate the bytes: white space, ':', ',', '-' and an optional "0x" in front of each byte. */
    private static final Pattern SEPARATORS = Pattern.compile("0[xX]|[\\s:,\\-]");

    private static final Pattern HEX_DIGITS = Pattern.compile("[0-9a-fA-F]*");

    private static final HexFormat HEX = HexFormat.of();

    private HexConvert() {
    }

    /**
     * String to hex string converter.
     *
     * @param textString text to convert
     * @param cs         text charset
     * @return lower case hexadecimal form of the encoded text, two digits per byte
     */
    public static String textStringToHexString(String textString, Charset cs) {
        Objects.requireNonNull(textString, "textString");
        Objects.requireNonNull(cs, "cs");
        return HEX.formatHex(textString.getBytes(cs));
    }

    /**
     * Hex string to string converter. Separators (spaces, new lines, {@code :}, {@code ,}, {@code -})
     * and {@code 0x} prefixes are allowed between the bytes.
     *
     * @param hexString hex to convert
     * @param cs        text charset
     * @return the decoded text
     * @throws IllegalArgumentException if the input has characters that are not hex digits, or an odd
     *                                  number of digits (the last byte would be incomplete)
     */
    public static String hexStringToTextString(String hexString, Charset cs) {
        Objects.requireNonNull(hexString, "hexString");
        Objects.requireNonNull(cs, "cs");
        String digits = normalize(hexString);
        return new String(HEX.parseHex(digits), cs);
    }

    /** Removes the separators and checks that only complete pairs of hex digits remain. */
    static String normalize(String hexString) {
        String digits = SEPARATORS.matcher(hexString).replaceAll("");
        if (!HEX_DIGITS.matcher(digits).matches()) {
            throw new IllegalArgumentException("Only hexadecimal digits (0-9, a-f) are allowed");
        }
        if (digits.length() % 2 != 0) {
            throw new IllegalArgumentException(
                    "The hex text has an odd number of digits (" + digits.length() + "); every byte needs two digits");
        }
        return digits;
    }
}
