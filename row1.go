package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func row1(outputEditor *widget.Entry) fyne.CanvasObject {
	outputEditor.SetPlaceHolder("Kết quả dịch từ AI sẽ xuất hiện ở đây...")
	outputEditor.Wrapping = fyne.TextWrapWord
	//rightPanelBg := canvas.NewRectangle(color.NRGBA{R: 255, G: 248, B: 235, A: 255})
	row1 := container.NewClip(outputEditor)
	return row1
}

type outputEditorTheme struct {
	fyne.Theme
}

func (t outputEditorTheme) Size(name fyne.ThemeSizeName) float32 {
	size := t.Theme.Size(name)
	if name == theme.SizeNameText {
		return size * 1.3
	}
	return size
}
