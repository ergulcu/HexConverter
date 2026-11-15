// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

/** The converter as a module, so that {@code javafx:jlink} can build a small runtime image. */
module com.ergulcu.hexconverter {
    requires javafx.controls;
    requires javafx.fxml;

    exports com.ergulcu.hex;
    exports com.ergulcu.hexconverter;

    opens com.ergulcu.hexconverter to javafx.fxml;
}
