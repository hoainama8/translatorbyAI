package main

import (
	"errors"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
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
	var popUp *widget.PopUp
	content := container.NewStack()
	popUp = widget.NewModalPopUp(container.NewThemeOverride(content, blueButtonTheme{Theme: theme.DefaultTheme()}), myWindow.Canvas())
	popUp.Resize(fyne.NewSize(1500, 700))

	historyBtn := widget.NewButton("Lịch sử dịch", func() {
		entries, err := readTranslationHistory(historyPath)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		closeButton := widget.NewButton("Đóng", func() {
			popUp.Hide()
		})
		closeButton.Importance = widget.HighImportance
		content.Objects = []fyne.CanvasObject{container.NewBorder(
			widget.NewLabelWithStyle("Thống kê lịch sử dịch tài liệu", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewHBox(layout.NewSpacer(), closeButton),
			nil,
			nil,
			container.NewThemeOverride(
				translationHistoryTableView(translationHistoryTable(entries)),
				outputEditorTheme{Theme: theme.DefaultTheme()},
			),
		)}
		content.Refresh()
		popUp.Show()
	})

	//Xem lịch sử sử dụng token

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
	historyBtn.Importance = widget.HighImportance
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

	res := container.NewVBox(container.NewThemeOverride(container.NewCenter(
		container.NewGridWithRows(3,
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
			container.NewGridWithColumns(2,
				//costBtn,
				historyBtn))),
		blueButtonTheme{Theme: theme.DefaultTheme()}))
	return res
}

func blueButton(text string, icon fyne.Resource, tapped func()) *widget.Button {
	button := widget.NewButtonWithIcon(text, icon, tapped)
	button.Importance = widget.HighImportance
	return button
}
