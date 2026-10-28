package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Tabmain(nametab string, text string, tabs *container.AppTabs) *widget.Entry {
	content := widget.NewMultiLineEntry()
	if text == "" {
		text = ("Noi dung dich sex xuat tai day")
	}
	content.SetText(text)
	content.Wrapping = fyne.TextWrapWord
	content.TextStyle = fyne.TextStyle{Bold: true}

	tabs.Append(container.NewTabItem(nametab, container.NewThemeOverride(
		container.NewClip(content),
		contentTabTheme{Theme: theme.DefaultTheme()},
	)))
	tabs.Refresh()
	return content
}

type contentTabTheme struct {
	fyne.Theme
}

func (t contentTabTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground && variant == theme.VariantLight {
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	}
	return t.Theme.Color(name, variant)
}

func (t contentTabTheme) Size(name fyne.ThemeSizeName) float32 {
	size := t.Theme.Size(name)
	if name == theme.SizeNameText {
		return size * 1.3
	}
	return size
}
