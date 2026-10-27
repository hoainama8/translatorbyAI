package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const donateDisplayDuration = 5 * time.Second

//go:embed icon.png
var sponsorImage []byte

// CỘT 1: INPUT
func createInputColumn(donations []Donate) (*fyne.Container, func()) {

	customLay := layout.CustomPaddedLayout{
		TopPadding:    20,
		BottomPadding: 20,
		LeftPadding:   30,
		RightPadding:  30,
	}
	currentDonate := Donate{}
	if len(donations) > 0 {
		donations = append([]Donate(nil), donations...)
		currentDonate = donations[0]
	}

	inputComponents := container.New(customLay, DonateContainer(currentDonate))
	refresh := func() {
		inputComponents.Objects = []fyne.CanvasObject{DonateContainer(currentDonate)}
		inputComponents.Refresh()
	}

	if len(donations) > 1 {
		currentIndex := 0
		animation := fyne.NewAnimation(donateDisplayDuration*time.Duration(len(donations)), func(progress float32) {
			nextIndex := int(progress * float32(len(donations)))
			if nextIndex >= len(donations) {
				nextIndex = len(donations) - 1
			}
			if nextIndex == currentIndex {
				return
			}

			currentIndex = nextIndex
			currentDonate = donations[currentIndex]
			refresh()
		})
		animation.Curve = fyne.AnimationLinear
		animation.RepeatCount = fyne.AnimationRepeatForever
		animation.Start()
	}

	return inputComponents, refresh
}

func DonateContainer(d Donate) fyne.CanvasObject {
	logoSize := fyne.NewSize(160, 160)
	logoData := d.Image
	if len(logoData) == 0 {
		logoData = sponsorImage
	}
	logo, err := circularLogo(logoData)
	var logoView fyne.CanvasObject
	if err != nil {
		logoView = widget.NewLabelWithStyle("Không thể hiển thị logo: "+err.Error(), fyne.TextAlignCenter, fyne.TextStyle{})
	} else {
		background := canvas.NewCircle(color.NRGBA{R: 42, G: 52, B: 72, A: 255})
		logoCanvas := canvas.NewImageFromImage(logo)
		logoCanvas.FillMode = canvas.ImageFillContain
		logoCanvas.SetMinSize(logoSize)
		border := canvas.NewCircle(color.Transparent)
		border.StrokeColor = color.NRGBA{R: 145, G: 205, B: 255, A: 255}
		border.StrokeWidth = 2
		logoView = container.NewGridWrap(logoSize, container.NewStack(background, logoCanvas, border))
	}

	details := container.NewVBox(
		widget.NewLabelWithStyle("NHÀ TÀI TRỢ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(d.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		donationDetails(d),
	)
	content := container.NewVBox(
		container.NewCenter(logoView),
		canvas.NewLine(color.NRGBA{R: 73, G: 88, B: 112, A: 255}),
		details,
	)
	darkBackground := canvas.NewRectangle(color.NRGBA{R: 24, G: 30, B: 43, A: 255})
	darkBackground.StrokeColor = color.NRGBA{R: 66, G: 81, B: 105, A: 255}
	darkBackground.StrokeWidth = 1

	return container.NewVScroll(container.NewStack(
		darkBackground,
		container.NewThemeOverride(
			container.NewPadded(content),
			donationTheme{Theme: theme.DefaultTheme()},
		),
	))
}

func donationDetails(d Donate) fyne.CanvasObject {
	details := make([]fyne.CanvasObject, 0, 4)
	for _, value := range []string{d.Address, d.Phone} {
		if value != "" {
			details = append(details, widget.NewLabel(value))
		}
	}
	if d.Website != "" {
		details = append(details, donationWebsite(d.Website))
	}
	if d.Logo != "" {
		details = append(details, widget.NewLabel(d.Logo))
	}
	if len(details) == 0 {
		details = append(details, widget.NewLabel("Thông tin liên hệ sẽ được cập nhật."))
	}
	return container.NewVBox(details...)
}

func donationWebsite(value string) fyne.CanvasObject {
	displayValue := strings.TrimSpace(value)
	parsedURL, err := url.Parse(displayValue)
	if err != nil || parsedURL.Scheme == "" {
		parsedURL, err = url.Parse("https://" + displayValue)
	}
	if err != nil || parsedURL.Host == "" || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return widget.NewLabel(displayValue)
	}
	return widget.NewHyperlink(displayValue, parsedURL)
}

type donationTheme struct {
	fyne.Theme
}

func (t donationTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground {
		return color.NRGBA{R: 239, G: 244, B: 255, A: 255}
	}
	if name == theme.ColorNameDisabled {
		return color.NRGBA{R: 160, G: 174, B: 196, A: 255}
	}
	return t.Theme.Color(name, variant)
}

func (t donationTheme) Size(name fyne.ThemeSizeName) float32 {
	size := t.Theme.Size(name)
	if name == theme.SizeNameText {
		return size * 1.25
	}
	return size
}

func circularLogo(data []byte) (image.Image, error) {
	source, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	bounds := source.Bounds()
	size := bounds.Dx()
	if bounds.Dy() < size {
		size = bounds.Dy()
	}
	if size == 0 {
		return nil, image.ErrFormat
	}

	offsetX := bounds.Min.X + (bounds.Dx()-size)/2
	offsetY := bounds.Min.Y + (bounds.Dy()-size)/2
	cropped := image.NewNRGBA(image.Rect(0, 0, size, size))
	draw.Draw(cropped, cropped.Bounds(), source, image.Pt(offsetX, offsetY), draw.Src)

	radius := float64(size) / 2
	center := radius
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			pixel := cropped.NRGBAAt(x, y)
			distance := math.Hypot(float64(x)+0.5-center, float64(y)+0.5-center)
			coverage := math.Max(0, math.Min(1, radius+0.5-distance))
			pixel.A = uint8(float64(pixel.A) * coverage)
			cropped.SetNRGBA(x, y, pixel)
		}
	}
	return cropped, nil
}
