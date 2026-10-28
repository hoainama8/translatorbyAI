package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func row1(outputText string, historyPath string, myWindow fyne.Window) (*container.AppTabs, *widget.Entry, func()) {
	tabs := container.NewAppTabs()

	outputEditor := Tabmain("Nội dung dịch", outputText, tabs)
	TabGuide(tabs)
	TabCopyright(tabs)
	refreshHistory := TabHistory(historyPath, myWindow, tabs)
	refreshCost := TabCost(myWindow, tabs)
	TabPrice(models, tabs)
	return tabs, outputEditor, func() {
		if refreshHistory != nil {
			refreshHistory()
		}
		if refreshCost != nil {
			refreshCost()
		}
	}
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
