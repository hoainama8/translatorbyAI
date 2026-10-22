package main

import (
	"errors"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func row2_2(inputDirectory string, outDirectory string, historyPath string, myApp fyne.App, myWindow fyne.Window, req *translationRequest, outputText *string) *fyne.Container {
	var popUp *widget.PopUp
	content := container.NewStack()
	popUp = widget.NewModalPopUp(content, myWindow.Canvas())
	popUp.Resize(fyne.NewSize(1500, 700))

	open := widget.NewButton("Lịch sử dịch", func() {
		entries, err := readTranslationHistory(historyPath)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		closeButton := widget.NewButton("Đóng", func() {
			popUp.Hide()
		})
		content.Objects = []fyne.CanvasObject{container.NewBorder(
			widget.NewLabelWithStyle("Thống kê lịch sử dịch tài liệu", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewHBox(layout.NewSpacer(), closeButton),
			nil,
			nil,
			container.NewThemeOverride(translationHistoryTable(entries), outputEditorTheme{Theme: theme.DefaultTheme()}),
		)}
		content.Refresh()
		popUp.Show()
	})

	// Xem lịch sử sử dụng token
	costBtn := widget.NewButton("Chi phí đã sử dụng", func() {
		popUp.Resize(fyne.NewSize(1300, 700))
		tokenUsagePath, err := defaultTokenUsagePath()
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		tokenFile, err := os.ReadFile(tokenUsagePath)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		tokenData, err := decodeTokenUsageData(tokenFile)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		tabs := container.NewAppTabs()
		for _, model := range models {
			table := CostTable(tokenData, model)
			for col, width := range []float32{210, 170, 220, 220, 220, 260} {
				table.SetColumnWidth(col, width)
			}
			tabs.Append(container.NewTabItemWithIcon(
				model,
				theme.DocumentIcon(),
				container.NewThemeOverride(table, outputEditorTheme{Theme: theme.DefaultTheme()}),
			))
		}

		closeButton := widget.NewButton("Đóng", func() {
			popUp.Hide()
		})
		content.Objects = []fyne.CanvasObject{container.NewBorder(
			widget.NewLabelWithStyle("Bảng thống kê chi phí", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewHBox(layout.NewSpacer(), closeButton),
			nil,
			nil,
			tabs,
		)}
		content.Refresh()
		popUp.Show()
	})
	pdfBtn := widget.NewButtonWithIcon("xuat ra file PDF", theme.FileTextIcon(), func() {
		if err := WritePdf(*outputText, *req, outDirectory, myWindow); err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		dialog.ShowInformation("Thành công", "Đã lưu tệp PDF vào thư mục tài liệu dịch.", myWindow)
	})
	wordBtn := widget.NewButtonWithIcon("Xuất ra file Word", theme.DocumentIcon(), func() {
		if *outputText == "" || outputText == nil {
			dialog.ShowError(errors.New("không có nội dung để lưu "), myWindow)
			return
		}
		path, err := saveTranslationFile(outDirectory, req.FileName, req.TargetLang, *outputText, ".docx")
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		dialog.ShowInformation("Thành công", "Đã lưu tệp Word: "+path, myWindow)
	})

	res := container.NewVBox(container.NewCenter(
		container.NewGridWithRows(2,
			container.NewGridWithColumns(2,
				widget.NewButtonWithIcon("Mở thư mục tài liệu gốc", theme.FolderOpenIcon(), OpenDir(inputDirectory, myApp, myWindow)),
				widget.NewButtonWithIcon("Mở thư mục tài liệu dịch", theme.FolderOpenIcon(), OpenDir(outDirectory, myApp, myWindow)),
			),
			container.NewGridWithColumns(4,
				open,
				costBtn,
				pdfBtn,
				wordBtn,
			))))
	return res
}

func translationHistoryTable(entries []translationHistoryEntry) *widget.Table {
	headers := []string{"Ngày", "Tên tài liệu", "Ngôn ngữ nguồn", "Ngôn ngữ đích", "Trạng thái"}
	rows := len(entries) + 1
	if len(entries) == 0 {
		rows++
	}
	table := widget.NewTable(
		func() (int, int) { return rows, len(headers) },
		func() fyne.CanvasObject {
			return container.NewStack(
				canvas.NewRectangle(theme.Color(theme.ColorNameBackground)),
				widget.NewLabel(""),
			)
		},
		func(cell widget.TableCellID, object fyne.CanvasObject) {
			cellContainer := object.(*fyne.Container)
			background := cellContainer.Objects[0].(*canvas.Rectangle)
			label := cellContainer.Objects[1].(*widget.Label)
			background.FillColor = theme.Color(theme.ColorNameBackground)
			label.Truncation = fyne.TextTruncateEllipsis
			if cell.Row == 0 {
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.SetText(headers[cell.Col])
				background.Refresh()
				return
			}
			label.TextStyle = fyne.TextStyle{}
			if cell.Row%2 == 1 {
				background.FillColor = theme.Color(theme.ColorNameInputBackground)
			}
			if len(entries) == 0 {
				if cell.Col == 0 {
					label.SetText("Chưa có lịch sử dịch")
				} else {
					label.SetText("")
				}
				background.Refresh()
				return
			}

			entry := entries[cell.Row-1]
			switch cell.Col {
			case 0:
				if entry.Date.IsZero() {
					label.SetText("")
				} else {
					label.SetText(entry.Date.Local().Format(time.DateTime))
				}
			case 1:
				label.SetText(entry.DocumentName)
			case 2:
				label.SetText(entry.DocumentLanguage)
			case 3:
				label.SetText(entry.TargetLanguage)
			case 4:
				label.SetText(entry.Result)
			default:
				label.SetText("")
			}
			background.Refresh()
		},
	)
	for col, width := range []float32{190, 420, 220, 220, 180} {
		table.SetColumnWidth(col, width)
	}
	return table
}
