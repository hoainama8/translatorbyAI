package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TabCost(myWindow fyne.Window, tabs *container.AppTabs) func() {

	tokenUsagePath, err := defaultTokenUsagePath()
	if err != nil {
		dialog.ShowError(err, myWindow)
		return nil
	}
	tokenFile, err := os.ReadFile(tokenUsagePath)
	if err != nil {
		dialog.ShowError(err, myWindow)
		return nil
	}
	tokenData, err := decodeTokenUsageData(tokenFile)
	if err != nil {
		dialog.ShowError(err, myWindow)
		return nil
	}

	tab := container.NewTabItem("Thống kê chi phí", costContent(tokenData))
	tabs.Append(tab)
	tabs.Refresh()
	return func() {
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
		tab.Content = costContent(tokenData)
		tabs.Refresh()
	}
}

func costContent(tokenData tokenUsageData) fyne.CanvasObject {
	tabsub := container.NewAppTabs()
	for _, model := range models {
		table := CostTable(tokenData, model)
		tabsub.Append(container.NewTabItem(model,
			container.NewThemeOverride(
				weightedTableView(table, []float32{120, 100, 110, 110, 110, 120}),
				outputEditorTheme{Theme: theme.DefaultTheme()},
			)))
	}
	tabsub.Refresh()
	return container.NewScroll(tabsub)
}

func CostTable(entries tokenUsageData, model string) *widget.Table {
	headers := []string{"Ngày", "Giờ", "Input", "Output", "Cached", "Chi phí"}
	sessions := entries[model]
	totalRow := len(sessions) + 1
	if len(sessions) == 0 {
		totalRow++
	}
	rows := totalRow + 1
	price := modelPrices[model]
	sessionCost := func(session tokenUsageSession) float64 {
		inputTokens := session.InputTokens - session.CachedTokens
		if inputTokens < 0 {
			inputTokens = 0
		}
		return float64(inputTokens)*price.Input +
			float64(session.CachedTokens)*price.Cached +
			float64(session.OutputTokens+session.ThoughtTokens)*price.Output
	}
	var totalCost float64
	for _, session := range sessions {
		totalCost += sessionCost(session)
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
			if cell.Row == totalRow {
				label.TextStyle = fyne.TextStyle{Bold: true}
				switch cell.Col {
				case 0:
					label.SetText("Tổng")
				case 5:
					label.SetText(fmt.Sprintf("$%.6f", totalCost))
				default:
					label.SetText("")
				}
				background.Refresh()
				return
			}
			if len(sessions) == 0 {
				if cell.Col == 0 {
					label.SetText("Chưa có dữ liệu")
				} else {
					label.SetText("")
				}
				background.Refresh()
				return
			}

			session := sessions[cell.Row-1]
			switch cell.Col {
			case 0:
				label.SetText(session.StartedAt.Format("2006-01-02"))
			case 1:
				label.SetText(session.StartedAt.Format("15:04:05"))
			case 2:
				label.SetText(fmt.Sprintf("%d", session.InputTokens))
			case 3:
				label.SetText(fmt.Sprintf("%d", session.OutputTokens+session.ThoughtTokens))
			case 4:
				label.SetText(fmt.Sprintf("%d", session.CachedTokens))
			case 5:
				label.SetText(fmt.Sprintf("$%.6f", sessionCost(session)))
			}
			background.Refresh()
		},
	)
	for col, width := range []float32{120, 100, 110, 110, 110, 120} {
		table.SetColumnWidth(col, width)
	}
	return table
}
