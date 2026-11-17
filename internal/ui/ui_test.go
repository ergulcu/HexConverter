// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

package ui

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func newView(t *testing.T) (*View, *[]error) {
	t.Helper()
	a := test.NewTempApp(t)
	w := a.NewWindow("test")
	t.Cleanup(w.Close)
	v := New(w)
	w.SetContent(v.Content)
	w.Resize(fyne.NewSize(900, 600))
	var errs []error
	v.ShowError = func(err error) { errs = append(errs, err) }
	return v, &errs
}

func TestTextToHexAndBack(t *testing.T) {
	v, errs := newView(t)

	test.Type(v.Input, "Ağaç")
	test.Tap(v.ToHex)
	if v.Output.Text != "41c49f61c3a7" {
		t.Errorf("hex = %q", v.Output.Text)
	}

	v.Input.SetText(v.Output.Text)
	test.Tap(v.ToText)
	if v.Output.Text != "Ağaç" || len(*errs) != 0 {
		t.Errorf("text = %q, errors = %v", v.Output.Text, *errs)
	}
}

func TestTheSelectedCharsetIsUsed(t *testing.T) {
	v, _ := newView(t)
	v.Charset.SetSelected("ISO-8859-9")
	v.Input.SetText("ğ")
	test.Tap(v.ToHex)
	if v.Output.Text != "f0" {
		t.Errorf("hex = %q", v.Output.Text)
	}
}

func TestInvalidHexShowsAnErrorAndKeepsTheResult(t *testing.T) {
	v, errs := newView(t)
	v.Output.SetText("previous")
	v.Input.SetText("48656")
	test.Tap(v.ToText)
	if len(*errs) != 1 || !strings.Contains((*errs)[0].Error(), "odd number of digits") {
		t.Errorf("errors = %v", *errs)
	}
	if v.Output.Text != "previous" {
		t.Errorf("output = %q", v.Output.Text)
	}
}

func TestCopyPutsTheResultOnTheClipboard(t *testing.T) {
	v, _ := newView(t)
	v.Input.SetText("Hi")
	test.Tap(v.ToHex)
	test.Tap(v.Copy)
	if got := fyne.CurrentApp().Clipboard().Content(); got != "4869" {
		t.Errorf("clipboard = %q", got)
	}
}

func TestUTF8IsSelectedFirst(t *testing.T) {
	v, _ := newView(t)
	if v.Charset.Selected != "UTF-8" || len(v.Charset.Options) != 6 {
		t.Errorf("selected %q of %v", v.Charset.Selected, v.Charset.Options)
	}
}
