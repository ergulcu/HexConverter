// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

// Command hexconverter is the desktop window: text on the left, hex on the right.
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/ergulcu/HexConverter/internal/ui"
)

func main() {
	a := app.NewWithID("com.ergulcu.hexconverter")
	w := a.NewWindow("Hex Converter")
	v := ui.New(w)
	w.SetContent(v.Content)
	w.Resize(fyne.NewSize(1024, 640))
	w.ShowAndRun()
}
