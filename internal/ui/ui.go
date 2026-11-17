// Copyright 2019 Mehmet ERGULCU
// SPDX-License-Identifier: Apache-2.0

// Package ui is the window of the converter: text on the left, the result on the right,
// the character set and the buttons in between.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/ergulcu/HexConverter/hexconv"
)

// View holds the widgets, so that tests can drive them.
type View struct {
	Input   *widget.Entry
	Output  *widget.Entry
	Charset *widget.Select
	ToText  *widget.Button
	ToHex   *widget.Button
	Copy    *widget.Button
	Content fyne.CanvasObject

	// ShowError reports a conversion error; a dialog on the window by default.
	ShowError func(error)
}

// New builds the view for the window.
func New(w fyne.Window) *View {
	v := &View{
		Input:  widget.NewMultiLineEntry(),
		Output: widget.NewMultiLineEntry(),
		ShowError: func(err error) {
			dialog.ShowError(err, w)
		},
	}
	for _, e := range []*widget.Entry{v.Input, v.Output} {
		e.Wrapping = fyne.TextWrapWord
		e.TextStyle = fyne.TextStyle{Monospace: true}
	}
	v.Input.SetPlaceHolder("Text or hex, e.g. 48 65 6c 6c 6f")
	v.Output.SetPlaceHolder("Result")

	v.Charset = widget.NewSelect(hexconv.Names(), nil)
	v.Charset.SetSelectedIndex(0)

	v.ToText = widget.NewButton("Hex To Text", func() {
		v.convert(hexconv.FromHex)
	})
	v.ToHex = widget.NewButton("Text To Hex", func() {
		v.convert(hexconv.ToHex)
	})
	v.Copy = widget.NewButton("Copy Result", func() {
		fyne.CurrentApp().Clipboard().SetContent(v.Output.Text)
	})

	middle := container.NewVBox(layout.NewSpacer(), v.Charset, v.ToText, v.ToHex, v.Copy, layout.NewSpacer())
	v.Content = container.New(&columns{}, v.Input, middle, v.Output)
	return v
}

func (v *View) convert(f func(string, hexconv.Charset) (string, error)) {
	cs, _ := hexconv.Lookup(v.Charset.Selected)
	out, err := f(v.Input.Text, cs)
	if err != nil {
		v.ShowError(err)
		return
	}
	v.Output.SetText(out)
}

// columns lays out three objects: the outer two share the width the middle one leaves.
type columns struct{}

func (*columns) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	pad := float32(8)
	mid := objects[1].MinSize().Width
	side := (size.Width - mid - 2*pad) / 2
	objects[0].Move(fyne.NewPos(0, 0))
	objects[0].Resize(fyne.NewSize(side, size.Height))
	objects[1].Move(fyne.NewPos(side+pad, 0))
	objects[1].Resize(fyne.NewSize(mid, size.Height))
	objects[2].Move(fyne.NewPos(side+mid+2*pad, 0))
	objects[2].Resize(fyne.NewSize(side, size.Height))
}

func (*columns) MinSize(objects []fyne.CanvasObject) fyne.Size {
	m := objects[1].MinSize()
	return fyne.NewSize(m.Width+2*200+16, m.Height)
}
