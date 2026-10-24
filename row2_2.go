package main

import (
	_ "embed"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//go:embed icon.png
var sponsorImage []byte

// CỘT 1: INPUT
func createInputColumn() (*fyne.Container, func()) {

	customLay := layout.CustomPaddedLayout{
		TopPadding:    20,
		BottomPadding: 20,
		LeftPadding:   30,
		RightPadding:  30,
	}
	inputComponents := container.New(customLay, DonateContainer())
	refresh := func() {
		inputComponents.Objects = []fyne.CanvasObject{DonateContainer()}
		inputComponents.Refresh()
	}

	return inputComponents, refresh
}

func DonateContainer() fyne.CanvasObject {
	image := canvas.NewImageFromResource(fyne.NewStaticResource("icon.png", sponsorImage))
	image.FillMode = canvas.ImageFillCover
	image.SetMinSize(fyne.NewSize(180, 100))
	donate := container.NewVScroll(widget.NewCard(
		"Thông tin nhà tài trợ",
		"Thông tin sẽ được cập nhật khi có nhà tài trợ.",
		container.NewGridWithColumns(2,
			image,
			container.NewVBox(
				widget.NewLabel("Tên: Chưa cập nhật"),
				widget.NewLabel("Địa chỉ liên hệ: Chưa cập nhật"),
				widget.NewLabel("Số điện thoại: Chưa cập nhật"),
				widget.NewLabel("Website: Chưa cập nhật"),
				widget.NewLabel("Logo: Chưa cập nhật"),
			),
		)))
	return donate
}
