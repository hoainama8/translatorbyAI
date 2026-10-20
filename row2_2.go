package main

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CỘT 1: INPUT
func createInputColumn(models []string) (*fyne.Container, func()) {

	tabs := createInputTabs(models)
	customLay := layout.CustomPaddedLayout{
		TopPadding:   20,
		LeftPadding:  30,
		RightPadding: 30,
	}
	inputComponents := container.New(customLay, tabs)
	refresh := func() {
		inputComponents.Objects = []fyne.CanvasObject{createInputTabs(models)}
		inputComponents.Refresh()
	}

	return inputComponents, refresh
}

func createInputTabs(models []string) *container.AppTabs {
	text := canvas.NewText("Ngày", color.Black)
	text.Resize(fyne.Size{Width: 1.5})
	// text.Color = color.Gray{Y: 80}
	tabs := container.NewAppTabs()
	var tokenData tokenUsageData
	tokenUsagePath, err := defaultTokenUsagePath()
	if err != nil {
		fmt.Println("Error finding token file:", err)
	} else {
		tokenFile, err := os.ReadFile(tokenUsagePath)
		if err != nil {
			fmt.Println("Error reading token file:", err)
		} else if tokenData, err = decodeTokenUsageData(tokenFile); err != nil {
			fmt.Println("Error parsing token file:", err)
		}
	}

	for _, model := range models {
		colDate := container.NewVBox(widget.NewLabelWithStyle(fmt.Sprintf("%s", text.Text), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		colHour := container.NewVBox(widget.NewLabelWithStyle("Giờ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		colInputToken := container.NewVBox(widget.NewLabelWithStyle("Input Token", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		colOuputToken := container.NewVBox(widget.NewLabelWithStyle("Output Token", fyne.TextAlignCenter, fyne.TextStyle{}))
		colThoughtToken := container.NewVBox(widget.NewLabelWithStyle("Thought Token", fyne.TextAlignCenter, fyne.TextStyle{}))
		colCachedToken := container.NewVBox(widget.NewLabelWithStyle("Cached Token", fyne.TextAlignCenter, fyne.TextStyle{}))
		colCost := container.NewVBox(widget.NewLabelWithStyle("Chi Phi", fyne.TextAlignCenter, fyne.TextStyle{}))

		price := modelPrices[model]
		var totalCost float64
		for _, session := range tokenData[model] {
			colDate.Add(widget.NewLabelWithStyle(session.StartedAt.Format("2006-01-02"), fyne.TextAlignCenter, fyne.TextStyle{}))
			colHour.Add(widget.NewLabelWithStyle(session.StartedAt.Format("15:04:05"), fyne.TextAlignCenter, fyne.TextStyle{}))
			colInputToken.Add(widget.NewLabelWithStyle(fmt.Sprintf("%d", session.InputTokens), fyne.TextAlignCenter, fyne.TextStyle{}))
			colOuputToken.Add(widget.NewLabelWithStyle(fmt.Sprintf("%d", session.OutputTokens), fyne.TextAlignCenter, fyne.TextStyle{}))
			colThoughtToken.Add(widget.NewLabelWithStyle(fmt.Sprintf("%d", session.ThoughtTokens), fyne.TextAlignCenter, fyne.TextStyle{}))
			colCachedToken.Add(widget.NewLabelWithStyle(fmt.Sprintf("%d", session.TotalTokens), fyne.TextAlignCenter, fyne.TextStyle{}))
			sessionCost := (float64(session.InputTokens) * price.Input) + (float64(session.OutputTokens+session.ThoughtTokens) * price.Output)
			totalCost += sessionCost
			colCost.Add(widget.NewLabelWithStyle(fmt.Sprintf("$%.6f", sessionCost), fyne.TextAlignCenter, fyne.TextStyle{}))
		}

		colDate.Add(widget.NewLabelWithStyle(text.ChangeText("Tổng").Text, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		colHour.Add(widget.NewLabel(""))
		colInputToken.Add(widget.NewLabel(""))
		colOuputToken.Add(widget.NewLabel(""))
		colThoughtToken.Add(widget.NewLabel(""))
		colCachedToken.Add(widget.NewLabel(""))
		colCost.Add(widget.NewLabelWithStyle(fmt.Sprintf("$%.6f", totalCost), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

		db := container.NewScroll(
			container.NewGridWithColumns(7, colDate, colHour, colInputToken,
				colOuputToken, colThoughtToken, colCachedToken, colCost))
		tabs.Append(container.NewTabItemWithIcon(
			model,
			theme.DocumentIcon(),
			container.NewStack(db),
		),
		)
	}

	priceModelColumn := container.NewVBox(widget.NewLabelWithStyle("Mô hình", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	priceInputColumn := container.NewVBox(widget.NewLabelWithStyle("Input (USD)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	priceOutputColumn := container.NewVBox(widget.NewLabelWithStyle("Output (USD)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	priceCachedColumn := container.NewVBox(widget.NewLabelWithStyle("Cached (USD)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	for _, model := range models {
		price := modelPrices[model]
		priceModelColumn.Add(widget.NewLabel(model))
		priceInputColumn.Add(widget.NewLabel(fmt.Sprintf("$%.6f", price.Input*limitToken)))
		priceOutputColumn.Add(widget.NewLabel(fmt.Sprintf("$%.6f", price.Output*limitToken)))
		priceCachedColumn.Add(widget.NewLabel(fmt.Sprintf("$%.6f", price.Cached*limitToken)))
	}

	priceTable := container.NewThemeOverride(container.NewGridWithColumns(4,
		priceModelColumn,
		priceInputColumn,
		priceOutputColumn,
		priceCachedColumn,
	), outputEditorTheme{Theme: theme.DefaultTheme()})
	tabs.Append(container.NewTabItemWithIcon(
		"Đơn giá",
		theme.DocumentIcon(),
		container.NewVBox(
			widget.NewLabelWithStyle("Đơn giá trên 1 triệu token", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			priceTable,
		),
	))

	return tabs
}
