// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

package com.ergulcu.hexconverter;

import com.ergulcu.hex.HexConvert;
import javafx.event.ActionEvent;
import javafx.fxml.FXML;
import javafx.scene.control.Alert;
import javafx.scene.control.ComboBox;
import javafx.scene.control.TextArea;

import java.nio.charset.Charset;
import java.util.List;

public class FXMLController {

    /** UTF-8 first; windows-1254 and ISO-8859-9 cover Turkish text written by older systems. */
    static final List<String> CHARSETS = List.of("UTF-8", "UTF-16", "ISO-8859-1", "ISO-8859-9", "windows-1254", "US-ASCII");

    @FXML
    private TextArea leftTextArea;

    @FXML
    private TextArea rightTextArea;

    @FXML
    private ComboBox<String> charsetComboBox;

    @FXML
    private void initialize() {
        charsetComboBox.getItems().setAll(CHARSETS.stream().filter(Charset::isSupported).toList());
        charsetComboBox.getSelectionModel().selectFirst();
    }

    @FXML
    private void handleHexToTextButtonAction(ActionEvent event) {
        try {
            rightTextArea.setText(HexConvert.hexStringToTextString(leftTextArea.getText(), selectedCharset()));
        } catch (IllegalArgumentException e) {
            showError("The text on the left is not valid hex", e.getMessage());
        }
    }

    @FXML
    private void handleTextToHexButtonAction(ActionEvent event) {
        rightTextArea.setText(HexConvert.textStringToHexString(leftTextArea.getText(), selectedCharset()));
    }

    private Charset selectedCharset() {
        return Charset.forName(charsetComboBox.getValue());
    }

    private static void showError(String header, String content) {
        Alert alert = new Alert(Alert.AlertType.ERROR);
        alert.setTitle("Hex Converter");
        alert.setHeaderText(header);
        alert.setContentText(content);
        alert.showAndWait();
    }
}
