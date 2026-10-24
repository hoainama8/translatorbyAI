package main

// func BtnCost(myWindow fyne.Window) widget.Button {
// 	content :=
// 	popUp := widget.NewPopUp(content)
// 	costBtn := widget.NewButton("Chi phí đã sử dụng", func() {
// 		popUp.Resize(fyne.NewSize(1300, 700))
// 		tokenUsagePath, err := defaultTokenUsagePath()
// 		if err != nil {
// 			dialog.ShowError(err, myWindow)
// 			return
// 		}
// 		tokenFile, err := os.ReadFile(tokenUsagePath)1
// 		if err != nil {
// 			dialog.ShowError(err, myWindow)
// 			return
// 		}
// 		tokenData, err := decodeTokenUsageData(tokenFile)
// 		if err != nil {
// 			dialog.ShowError(err, myWindow)
// 			return
// 		}

// 		tabs := container.NewAppTabs()
// 		for _, model := range models {
// 			table := CostTable(tokenData, model)
// 			for col, width := range []float32{210, 170, 220, 220, 220, 260} {
// 				table.SetColumnWidth(col, width)
// 			}
// 			tabs.Append(container.NewTabItemWithIcon(
// 				model,
// 				theme.DocumentIcon(),
// 				container.NewThemeOverride(table, outputEditorTheme{Theme: theme.DefaultTheme()}),
// 			))
// 		}

// 		closeButton := widget.NewButton("Đóng", func() {
// 			popUp.Hide()
// 		})
// 		closeButton.Importance = widget.HighImportance
// 		content.Objects = []fyne.CanvasObject{container.NewBorder(
// 			widget.NewLabelWithStyle("Bảng thống kê chi phí", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
// 			container.NewHBox(layout.NewSpacer(), closeButton),
// 			nil,
// 			nil,
// 			tabs,
// 		)}
// 		content.Refresh()
// 		popUp.Show()
// 	})
// 	return *costBtn
// }
