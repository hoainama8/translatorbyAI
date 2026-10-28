package main

import (
	"errors"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type blueButtonTheme struct {
	fyne.Theme
}

func (t blueButtonTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return theme.PrimaryColorNamed(theme.ColorBrown)
	case theme.ColorNameForegroundOnPrimary:
		return color.White
	default:
		return t.Theme.Color(name, variant)
	}
}

func row2_2(inputDirectory string, outDirectory string, historyPath string, myApp fyne.App, myWindow fyne.Window, req *translationRequest, outputText *string) *fyne.Container {

	pdfBtn := widget.NewButtonWithIcon("xuat ra file PDF", theme.FileTextIcon(), func() {
		if outputText == nil || *outputText == "" {
			dialog.ShowError(errors.New("không có nội dung để lưu "), myWindow)
			return
		}
		if err := WritePdf(*outputText, *req, outDirectory, myWindow); err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		dialog.ShowInformation("Thành công", "Đã lưu tệp PDF vào thư mục tài liệu dịch.", myWindow)
	})
	//historyBtn.Importance = widget.HighImportance
	//costBtn.Importance = widget.HighImportance
	pdfBtn.Importance = widget.HighImportance
	exportButton := func(label, extension string) *widget.Button {
		return widget.NewButtonWithIcon(label, theme.DocumentIcon(), func() {
			if outputText == nil || *outputText == "" {
				dialog.ShowError(errors.New("không có nội dung để lưu "), myWindow)
				return
			}
			path, err := saveTranslationFile(outDirectory, req.FileName, req.TargetLang, *outputText, extension)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			dialog.ShowInformation("Thành công", "Đã lưu tệp: "+path, myWindow)
		})
	}
	wordBtn := exportButton("Xuất ra file Word", ".docx")
	txtBtn := exportButton("Xuất ra file TXT", ".txt")
	odtBtn := exportButton("Xuất ra file ODT", ".odt")
	wordBtn.Importance = widget.HighImportance
	txtBtn.Importance = widget.HighImportance
	odtBtn.Importance = widget.HighImportance

	res := container.NewVBox(container.NewThemeOverride(
		container.NewCenter(container.NewGridWithRows(2,
			container.NewGridWithColumns(2,
				blueButton("Mở thư mục tài liệu gốc", theme.FolderOpenIcon(), OpenDir(inputDirectory, myApp, myWindow)),
				blueButton("Mở thư mục tài liệu dịch", theme.FolderOpenIcon(), OpenDir(outDirectory, myApp, myWindow)),
			),
			container.NewGridWithColumns(4,
				pdfBtn,
				wordBtn,
				txtBtn,
				odtBtn,
			),
		)),
		blueButtonTheme{Theme: theme.DefaultTheme()},
	))
	return res
}

func blueButton(text string, icon fyne.Resource, tapped func()) *widget.Button {
	button := widget.NewButtonWithIcon(text, icon, tapped)
	button.Importance = widget.HighImportance
	return button
}
