package com.ergulcu.hex;

import java.nio.charset.Charset;

/**
 *
 * @author Mehmet ERGULCU
 */
public class HexConvert {

    /**
     * String to Hex String Converter
     * @param textString text to convert
     * @param cs text charset
     * @return hexadecimal equivalence of text
     */
    public static String textStringToHexString(String textString, Charset cs) {
        byte[] hashInBytes=textString.getBytes(cs);
        StringBuilder sb = new StringBuilder();
        for (byte b : hashInBytes) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }

    /**
     * Hex String to String Converter
     * @param hexString hex to convert
     * @param cs text charset
     * @return text equivalence of hexadecimal text
     */
    public static String hexStringToTextString(String hexString, Charset cs) {

        byte[] bytes = new byte[hexString.length() / 2];
        for (int i = 0; i < hexString.length() - 1; i += 2) {
            //get the hex in pairs
            String output = hexString.substring(i, (i + 2));
            //convert hex to decimal
            int decimal = Integer.parseInt(output, 16);
            //convert the decimal to character
            bytes[i / 2] = (byte) decimal;
        }

        return new String(bytes, cs);
    }

}
