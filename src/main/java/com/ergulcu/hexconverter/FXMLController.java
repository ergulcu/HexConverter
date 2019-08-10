package com.ergulcu.hexconverter;

import com.ergulcu.hex.HexConvert;
import java.net.URL;
import java.nio.charset.Charset;
import java.util.ResourceBundle;
import javafx.event.ActionEvent;
import javafx.fxml.FXML;
import javafx.fxml.Initializable;
import javafx.scene.control.ComboBox;
import javafx.scene.control.TextArea;

/**
 *
 * @author Mehmet ERGULCU
 */
public class FXMLController implements Initializable {

    @FXML
    TextArea leftTextArea;

    @FXML
    TextArea rightTextArea;

    @FXML
    ComboBox<String> charsetComboBox;

    @Override
    public void initialize(URL url, ResourceBundle rb) {
        charsetComboBox.getItems().addAll("UTF-8", "ISO-8859-1", "ISO-8859-9", "US-ASCII");
        charsetComboBox.getSelectionModel().selectFirst();
    }

    @FXML
    private void handleHexToTextButtonAction(ActionEvent event) {
        String leftText = getLeftText();
        leftText = cleanHex(leftText);
        Charset selectedCharset = getSelectedCharset();
        String rightText = HexConvert.hexStringToTextString(leftText, selectedCharset);
        setRightText(rightText);
    }

    @FXML
    private void handleTextToHexButtonAction(ActionEvent event) {
        String leftText = getLeftText();
        Charset selectedCharset = getSelectedCharset();
        String rightText = HexConvert.textStringToHexString(leftText, selectedCharset);
        setRightText(rightText);
    }

    private String getLeftText() {
        String text = leftTextArea.getText();
        return text;
    }

    private void setRightText(String rightText) {
        rightTextArea.setText(rightText);
    }

    private String cleanHex(String text) {
        return text.replaceAll("[^a-fA-F0-9]", "");
    }

    private Charset getSelectedCharset() {
        return Charset.forName(charsetComboBox.getValue());
    }
}
