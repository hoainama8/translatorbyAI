package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func row1(outputText string, historyPath string, myWindow fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs()

	Tabmain("Nội dung dịch", outputText, tabs)
	TabGuide(tabs)
	TabCopyright(tabs)
	TabHistory(historyPath, myWindow, tabs)
	TabCost(myWindow, tabs)
	TabPrice(models, tabs)
	return tabs
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
