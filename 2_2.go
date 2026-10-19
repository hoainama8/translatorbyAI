package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func row2_2(inputDirectory string, outDirectory string, myApp fyne.App, myWindow fyne.Window) *fyne.Container {
	res := container.NewVBox(container.NewCenter(
		container.NewHBox(
			widget.NewButtonWithIcon("mo thu muc input", theme.FolderOpenIcon(), OpenDir(inputDirectory, myApp, myWindow)),
			widget.NewButtonWithIcon("mo thu muc output", theme.FolderOpenIcon(), OpenDir(outDirectory, myApp, myWindow)),
		)))
	return res
}
