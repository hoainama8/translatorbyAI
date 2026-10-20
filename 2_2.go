package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func row2_2(inputDirectory string, outDirectory string, myApp fyne.App, myWindow fyne.Window) *fyne.Container {
	var popUp *widget.PopUp
	content := container.NewVBox(
		widget.NewLabel("Content of the pop up"),
		widget.NewButton("Close", func() {
			popUp.Hide()
		}),
	)
	popUp = widget.NewModalPopUp(content, myWindow.Canvas())

	open := widget.NewButton("Show pop up", func() {
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
