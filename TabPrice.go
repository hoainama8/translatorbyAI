package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var priceTableColumnWidths = [...]float32{180, 110, 110, 110}

func TabPrice(models []string, tabs *container.AppTabs) {
	headers := []string{"Mô hình", "Input (USD)", "Output (USD)", "Cached (USD)"}
	priceTable := widget.NewTable(
		func() (int, int) {
			return len(models) + 1, len(headers)
		},
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
			model := models[cell.Row-1]
			price := modelPrices[model]
			switch cell.Col {
			case 0:
				label.SetText(model)
			case 1:
				label.SetText(fmt.Sprintf("$%.6f", price.Input*limitToken))
			case 2:
				label.SetText(fmt.Sprintf("$%.6f", price.Output*limitToken))
			case 3:
				label.SetText(fmt.Sprintf("$%.6f", price.Cached*limitToken))
			}
			background.Refresh()
		},
	)
	for col, width := range priceTableColumnWidths {
		priceTable.SetColumnWidth(col, width)
	}
	tabs.Append(container.NewTabItem(
		"Đơn giá",
		container.NewThemeOverride(
			weightedTableView(priceTable, priceTableColumnWidths[:]),
			outputEditorTheme{Theme: theme.DefaultTheme()},
		),
	))
	tabs.Refresh()
}
