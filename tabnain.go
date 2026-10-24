package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Tabmain(nametab string, text string, tabs *container.AppTabs) {
	content := widget.NewMultiLineEntry()
	if text == "" {
		text = ("Noi dung dich sex xuat tai day")
	}
	content.SetText(text)
	content.Wrapping = fyne.TextWrapWord

	tabs.Append(container.NewTabItem(nametab, container.NewThemeOverride(
		container.NewClip(content),
		outputEditorTheme{Theme: theme.DefaultTheme()},
	)))
	tabs.Refresh()
}
