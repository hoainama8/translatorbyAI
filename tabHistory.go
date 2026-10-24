package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TabHistory(historyPath string, mywindow fyne.Window, tabs *container.AppTabs) {
	entries, err := readTranslationHistory(historyPath)
	if err != nil {
		dialog.ShowError(err, mywindow)
	}
	tabs.Append(container.NewTabItem("Lich su dich tai lieu",
		container.NewThemeOverride(
			translationHistoryTableView(translationHistoryTable(entries)),
			outputEditorTheme{Theme: theme.DefaultTheme()},
		)))
	tabs.Refresh()
}

var translationHistoryColumnWidths = [...]float32{190, 420, 220, 220, 180}

type weightedTableLayout struct {
	table       *widget.Table
	columnWidth []float32
}

func (l weightedTableLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) == 0 {
		return
	}

	totalWidth := float32(0)
	for _, width := range l.columnWidth {
		totalWidth += width
	}
	availableWidth := size.Width - float32(len(l.columnWidth)-1)*l.table.Theme().Size(theme.SizeNamePadding)
	if size.Width > 0 && totalWidth > 0 {
		if availableWidth < 0 {
			availableWidth = 0
		}
		for col, width := range l.columnWidth {
			l.table.SetColumnWidth(col, availableWidth*width/totalWidth)
		}
	}
	objects[0].Move(fyne.NewPos(0, 0))
	objects[0].Resize(size)
}

func (l weightedTableLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}
	return objects[0].MinSize()
}

func translationHistoryTableView(table *widget.Table) *fyne.Container {
	return weightedTableView(table, translationHistoryColumnWidths[:])
}

func weightedTableView(table *widget.Table, columnWidth []float32) *fyne.Container {
	return container.New(weightedTableLayout{table: table, columnWidth: columnWidth}, table)
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
	for col, width := range translationHistoryColumnWidths {
		table.SetColumnWidth(col, width)
	}
	return table
}
