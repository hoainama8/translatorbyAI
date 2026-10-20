package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func row2_2(inputDirectory string, outDirectory string, historyPath string, myApp fyne.App, myWindow fyne.Window) *fyne.Container {
	var popUp *widget.PopUp
	content := container.NewVBox()
	popUp = widget.NewModalPopUp(content, myWindow.Canvas())
	popUp.Resize(fyne.NewSize(1000, 700))

	open := widget.NewButton("Lịch sử dịch", func() {
		entries, err := readTranslationHistory(historyPath)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		closeButton := widget.NewButton("Đóng", func() {
			popUp.Hide()
		})
		content.Objects = []fyne.CanvasObject{container.NewGridWithRows(2,
			translationHistoryTable(entries),
			container.NewCenter(container.NewGridWrap(fyne.NewSize(80, 20), closeButton)),
		)}
		content.Refresh()
		popUp.Show()
	})

	res := container.NewVBox(container.NewCenter(
		container.NewHBox(
			widget.NewButtonWithIcon("Mở thư mục tài liệu gốc", theme.FolderOpenIcon(), OpenDir(inputDirectory, myApp, myWindow)),
			widget.NewButtonWithIcon("Mở thư mục tài liệu dịch", theme.FolderOpenIcon(), OpenDir(outDirectory, myApp, myWindow)),
			open,
		)))
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
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(cell widget.TableCellID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			if cell.Row == 0 {
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.SetText(headers[cell.Col])
				return
			}
			label.TextStyle = fyne.TextStyle{}
			if len(entries) == 0 {
				if cell.Col == 1 {
					label.SetText("Chưa có lịch sử dịch")
				} else {
					label.SetText("")
				}
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
		},
	)
	for col, width := range []float32{150, 220, 150, 150, 110} {
		table.SetColumnWidth(col, width)
	}
	//table.UpdateHeader()
	return table
}
