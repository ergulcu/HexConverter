// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

package com.ergulcu.hex;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class HexConvertTest {

    private static final Charset TURKISH_ISO = Charset.forName("ISO-8859-9");

    @Test
    void textToHexInUtf8() {
        assertEquals("48656c6c6f", HexConvert.textStringToHexString("Hello", StandardCharsets.UTF_8));
        // ğ is two bytes in UTF-8
        assertEquals("c49f", HexConvert.textStringToHexString("ğ", StandardCharsets.UTF_8));
    }

    @Test
    void theCharsetChangesTheBytes() {
        assertEquals("f0", HexConvert.textStringToHexString("ğ", TURKISH_ISO));
        assertEquals("ğ", HexConvert.hexStringToTextString("f0", TURKISH_ISO));
    }

    @ParameterizedTest
    @ValueSource(strings = {"Hello", "Şişli'de ılık bir gün", "", "line 1\nline 2", "emoji 🙂"})
    void roundTripsInUtf8(String text) {
        String hex = HexConvert.textStringToHexString(text, StandardCharsets.UTF_8);
        assertEquals(text, HexConvert.hexStringToTextString(hex, StandardCharsets.UTF_8));
    }

    @ParameterizedTest
    @ValueSource(strings = {"48656C6C6F", "48 65 6c 6c 6f", "48:65:6c:6c:6f", "0x48 0x65 0x6C 0x6C 0x6F", "48-65-6c\n6c-6f", "48,65,6c,6c,6f"})
    void acceptsCommonSeparatorsAndUpperCase(String hex) {
        assertEquals("Hello", HexConvert.hexStringToTextString(hex, StandardCharsets.UTF_8));
    }

    @Test
    void rejectsAnOddNumberOfDigits() {
        // The old version silently dropped the last digit and returned "He".
        IllegalArgumentException error = assertThrows(IllegalArgumentException.class,
                () -> HexConvert.hexStringToTextString("48656", StandardCharsets.UTF_8));
        assertTrue(error.getMessage().contains("odd number of digits"));
    }

    @Test
    void rejectsCharactersThatAreNotHexDigits() {
        // The old version silently removed them, so "hello" was read as the single digit "e".
        assertThrows(IllegalArgumentException.class,
                () -> HexConvert.hexStringToTextString("hello", StandardCharsets.UTF_8));
        assertThrows(IllegalArgumentException.class,
                () -> HexConvert.hexStringToTextString("48zz", StandardCharsets.UTF_8));
    }

    @Test
    void emptyInputGivesEmptyOutput() {
        assertEquals("", HexConvert.hexStringToTextString("  \n ", StandardCharsets.UTF_8));
    }
}
